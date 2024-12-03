package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/term"
)

func readFile(filename string, gcm cipher.AEAD) (string, error) {
	ciphertext, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("Error reading file: %w", err)
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
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: crypter.exe [--read] <filename>")
		return
	}
	var filename string
	readMode := false
	if os.Args[1] == "--read" || os.Args[1] == "-r" {
		if len(os.Args) < 3 {
			fmt.Println("Error: Specify the filename to read.")
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
			fmt.Println("Error: Filename contains invalid characters. Avoid using / \\ : * ? \" < > | and control characters.")
			return
		}
	}
	key := []byte("examplekey123456examplekey123456") // Replace with secure key generation
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
		fmt.Print("\033[2J\033[H")
		cont, err := readFile(filename, gcm)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(cont)
		return
	}
	var ogcont string
	if _, err := os.Stat(filename); err == nil {
		ogcont, err = readFile(filename, gcm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Print("\033[2J\033[H")
	var textBuffer strings.Builder
	if ogcont != "" {
		textBuffer.WriteString(ogcont)
	}
	fmt.Println("Simple Text Editor (Press Ctrl+Q to save and exit)")
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
			fmt.Print("\033[2J\033[H")
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
	plaintext := []byte(textBuffer.String())
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("Error generating nonce:", err)
		return
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	err = os.WriteFile(filename, ciphertext, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Printf("Encrypted text saved to %s\n", filename)
}
