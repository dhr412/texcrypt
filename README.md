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
