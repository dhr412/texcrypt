# Crypter - Secure File Encryption with AES

Crypter is a command line text editor and file encryption utility that allows the user to securely encrypt and decrypt file contents using AES (Advanced Encryption Standard). The program can generate a secure AES key, store it with the file in an encrypted form, and retrieve the data when needed, ensuring data confidentiality.

## Purpose

Crypter allows the user to securely store data in a way that prevents exposure in case of a breach. The program enables:
- Encrypting and decrypting file contents.
- Simple text editing functionality, where users can modify file contents and save it securely.

## Prerequisites

- Go 1.16 or higher
- Access to a terminal or command prompt

## Installation

1. Clone or download the repository to your local machine.

    ```bash
    git clone https://github.com/DCoder206/crypt-text.git
    cd src
    ```

2. Build the Go program.

    ```bash
    go build -o crypter
    ```

3. The program is now ready to use.

## Usage

### Encrypting a file
1. **Run the program** with the file name and provide a passphrase (e.g., `securepassphrase`):
   
   ```bash
    ./crypter myfile
   ```

2. The program will generate an AES key, encrypt the file contents, and store the key encrypted alongside the file.

3. You will be prompted with a simple text editor to modify the content interactively. Press **Ctrl+Q** to save and exit.

### Decrypting a File
1. **Run the program in read mode** and provide the passphrase used during encryption:
   
    ```bash
    ./crypter --read myfile
    ```

2. The program will extract and decrypt the AES key from the file using the provided passphrase and then decrypt the file contents.

3. The decrypted contents of the file will be displayed
