package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: crypter.exe <filename>")
		return
	}
	filename := os.Args[1]
	if filepath.Ext(filename) != "" {
		fmt.Println("Error: Filename must not contain an extension. Use a plain name (e.g. 'myfile')")
		return
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
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("Error generating nonce:", err)
		return
	}
	fmt.Println("Enter text (Press Ctrl+Q to save):")
	reader := bufio.NewReader(os.Stdin)
	var input strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		if strings.TrimSpace(line) == "\x11" { // Ctrl+Q ASCII code
			break
		}
		input.WriteString(line)
	}
	plaintext := []byte(input.String())
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	err = os.WriteFile(filename, ciphertext, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Printf("Encrypted text saved to %s\n", filename)
}
