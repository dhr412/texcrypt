package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/term"
)

func generateSalt(size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("Invalid salt size: must be greater than 0")
	}
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("Error generating salt: %w", err)
	}
	return salt, nil
}
func generateEncKey(passphrase, salt []byte, keySize int) ([]byte, error) {
	if len(salt) == 0 {
		return nil, fmt.Errorf("Salt is required for key derivation")
	}
	return pbkdf2.Key(passphrase, salt, 10000, keySize, sha256.New), nil
}
func readFile(filename, passphrase string, saltSize int) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) || len(data) == 0 {
			return "", nil
		}
		return "", fmt.Errorf("Error reading file: %w", err)
	}
	if len(data) < saltSize {
		fmt.Println("File too small to be encrypted")
		return "", nil
	}
	salt, ciphertext := data[:saltSize], data[saltSize:]
	key, err := generateEncKey([]byte(passphrase), salt, 32)
	if err != nil {
		return "", fmt.Errorf("Error deriving key: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("Error creating cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("Error creating GCM: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("Error: Corrupted file or invalid content")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("Error decrypting file: %w", err)
	}
	return string(plaintext), nil
}
func clearScrn() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: crypter.exe [--read] <filename>")
		return
	}
	var filename string
	readMode := false
	if os.Args[1] == "--read" || os.Args[1] == "-r" {
		if len(os.Args) < 3 {
			fmt.Println("Error: Specify the filename to read")
			return
		}
		readMode = true
		filename = os.Args[2]
	} else {
		filename = os.Args[1]
	}
	if filepath.Ext(filename) != "" {
		fmt.Println("Error: Filename must not contain an extension. Use a plain name (e.g. 'myfile')")
		return
	}
	invalidChars := []rune{'/', '\\', ':', '*', '?', '"', '<', '>', '|'}
	for _, char := range filename {
		if unicode.IsControl(char) || strings.ContainsRune(string(invalidChars), char) {
			fmt.Println("Error: Filename contains invalid characters. Avoid using / \\ : * ? \" < > | and control characters")
			return
		}
	}
	fmt.Print("Enter passphrase: ")
	passphr, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading passphrase:", err)
		return
	}
	fmt.Println()
	salt, err := generateSalt(32)
	if err != nil {
		fmt.Println("Error generating salt:", err)
	}
	key, err := generateEncKey([]byte(passphr), salt, 32)
	if err != nil {
		fmt.Println("Error deriving key:", err)
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating GCM:", err)
		return
	}
	if readMode {
		clearScrn()
		cont, err := readFile(filename, string(passphr), 32)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(cont)
		return
	}
	var ogcont string
	if _, err := os.Stat(filename); err == nil {
		ogcont, err = readFile(filename, string(passphr), 32)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("Creating a new file...")
	}
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	clearScrn()
	var textBuffer strings.Builder
	if ogcont != "" {
		textBuffer.WriteString(ogcont)
	}
	fmt.Println("Text Editor (Press Ctrl+Q to save and exit)")
	fmt.Println("-------------------------------------------------")
	fmt.Print(textBuffer.String())
	for {
		buf := make([]byte, 1)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		if buf[0] == 17 {
			clearScrn()
			break
		}
		if buf[0] == 8 {
			if textBuffer.Len() > 0 {
				currentText := textBuffer.String()
				lastChar := currentText[len(currentText)-1]
				currentText = currentText[:len(currentText)-1]
				textBuffer.Reset()
				textBuffer.WriteString(currentText)
				if lastChar == '\n' {
					fmt.Print("\033[A\033[K")
				} else {
					fmt.Print("\b \b")
				}
				fmt.Print("\033[2K\r")
				fmt.Print(currentText[strings.LastIndex(currentText, "\n")+1:])
			} else {
				continue
			}
			continue
		}
		if buf[0] == '\r' || buf[0] == '\n' {
			textBuffer.WriteByte('\n')
			fmt.Print("\n")
			continue
		}
		textBuffer.WriteByte(buf[0])
		fmt.Print(string(buf[0]))
	}
	if textBuffer.Len() == 0 {
		fmt.Println("No changes made. File remains unaltered.")
		return
	}
	plaintext := []byte(textBuffer.String())
	if len(plaintext) == 0 {
		fmt.Println("No content to save. File will remain empty.")
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("Error generating nonce:", err)
		return
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	metaData := append(salt, nonce...)
	finalData := append(metaData, ciphertext...)
	err = os.WriteFile(filename, finalData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Printf("Encrypted text saved to %s\n", filename)
}
