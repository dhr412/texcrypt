package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/argon2"
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
	return argon2.IDKey(passphrase, salt, 2, 32*1024, 2, uint32(keySize)), nil
}

func encryptFile(filePath string, password []byte) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	salt, err := generateSalt(32)
	if err != nil {
		return err
	}

	key, err := generateEncKey(password, salt, 32)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	finalData := append(salt, append(nonce, ciphertext...)...)

	outFile := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".encrypt"
	return os.WriteFile(outFile, finalData, 0644)
}

func decryptFile(filePath string, password []byte) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}
	if len(data) < 32 {
		return fmt.Errorf("invalid encrypted file")
	}

	salt := data[:32]
	key, err := generateEncKey(password, salt, 32)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < 32+nonceSize {
		return fmt.Errorf("invalid encrypted content")
	}

	nonce := data[32 : 32+nonceSize]
	ciphertext := data[32+nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	outFile := strings.TrimSuffix(filePath, ".encrypt") + "_decrypted.txt"
	return os.WriteFile(outFile, plaintext, 0644)
}

func main() {
	encryptPath := flag.String("encrypt", "", "Encrypt the specified .txt or .md file")
	decryptPath := flag.String("decrypt", "", "Decrypt the specified .encrypt file")
	helpFlag := flag.Bool("help", false, "Show usage information")

	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  --encrypt=<filename.txt|filename.md>")
		fmt.Println("  --decrypt=<filename.encrypt>")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if *encryptPath != "" && *decryptPath != "" {
		fmt.Println("Error: Cannot use both --encrypt and --decrypt at the same time.")
		return
	}

	if *encryptPath != "" {
		if ext := filepath.Ext(*encryptPath); ext != ".txt" && ext != ".md" {
			fmt.Println("Error: Only .txt or .md files can be encrypted.")
			return
		}

		fmt.Print("Enter password: ")
		pass1, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		fmt.Print("Re-enter password: ")
		pass2, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		if string(pass1) != string(pass2) {
			fmt.Println("Error: Passwords do not match.")
			return
		}

		if err := encryptFile(*encryptPath, pass1); err != nil {
			fmt.Println("Encryption failed:", err)
		} else {
			fmt.Println("File encrypted successfully.")
		}
		return
	}

	if *decryptPath != "" {
		if filepath.Ext(*decryptPath) != ".encrypt" {
			fmt.Println("Error: Only files with .encrypt extension can be decrypted.")
			return
		}

		fmt.Print("Enter password: ")
		pass, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		if err := decryptFile(*decryptPath, pass); err != nil {
			fmt.Println("Decryption failed:", err)
		} else {
			fmt.Println("File decrypted successfully.")
		}
		return
	}

	fmt.Println("Error: You must specify either --encrypt or --decrypt.")
	flag.Usage()
}
