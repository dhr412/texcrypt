TexCrypt

`TexCrypt` is a lightweight, secure, and user-friendly CLI tool built in Go for encrypting and decrypting text-based files using AES-256-GCM. Designed with simplicity and data confidentiality in mind, TexCrypt ensures that sensitive file contents can be safely encrypted using strong passwords and modern cryptographic standards.

Whether you're storing personal notes, secure configuration, or secrets, `TexCrypt` provides a seamless way to lock down your data with minimal friction.

---

## Features

* **AES-256-GCM Encryption** – Provides both confidentiality and integrity.
* **Password-Based Encryption** – Uses PBKDF2 + SHA-256 to derive secure keys from passwords.
* **Salt + Nonce Randomization** – Each encryption operation uses unique salt and nonce values.
* **Simple CLI Interface** – Easy to use with only a couple of flags.
* **Secure File Output** – Encrypted output stored as `.encrypt`; decrypted output written as `_decrypted.txt`.
* **Safe Password Entry** – Prompts for password without echoing to the terminal.

---

## Installation

### From Prebuilt Releases

1. Visit the [Releases](https://github.com/dhr412/texcrypt/releases) page.
2. Download the binary for your platform.
3. Make it executable:

   ```bash
   chmod +x texcrypt
   ```

4. Run it:

   ```bash
   ./texcrypt --help
   ```

### Compiling from Source

Ensure you have [Go 1.20+](https://golang.org/dl/) installed:

```bash
git clone https://github.com/dhr412/texcrypt.git
cd texcrypt/src
go build -o texcrypt
./texcrypt --help
```

---

## Usage

```bash
texcrypt --encrypt=<file> | --decrypt=<file> [--help]
```

### Flags

* `--encrypt=<file>` – Encrypt a `.txt` or `.md` file using a password.
* `--decrypt=<file>` – Decrypt a `.encrypt` file using the original password.
* `--help` – Show usage instructions.

> Only `.txt` and `.md` files are allowed for encryption. Only `.encrypt` files are valid for decryption.

---

### Example

#### Encrypting a File

```bash
texcrypt --encrypt=secrets.txt
```

* Prompts for password and confirmation.
* Outputs: `secrets.encrypt`

#### Decrypting a File

```bash
texcrypt --decrypt=secrets.encrypt
```

* Prompts for the same password used to encrypt.
* Outputs: `secrets_decrypted.txt`

---

## How It Works

1. **Input Handling**:

   * Flags are parsed using Go’s `flag` package.
   * Validates file types and required options.

2. **Encryption**:

   * Prompts for password twice (with hidden input).
   * Generates 32-byte random salt and nonce.
   * Derives a 256-bit AES key using PBKDF2 (10000 iterations).
   * Encrypts file content using AES-GCM.
   * Writes `[salt][nonce][ciphertext]` to a `.encrypt` file.

3. **Decryption**:

   * Prompts for password (once).
   * Extracts salt and nonce from the file.
   * Re-derives the key using PBKDF2.
   * Decrypts the data and writes plaintext to `<original>_decrypted.txt`.

4. **Security Measures**:

   * Password input is hidden.
   * AES-GCM provides both encryption and authentication.
   * Key derivation is hardened using salt and PBKDF2.

---

## License

This project is open-sourced under the MIT license. Contributions, forks, and suggestions are welcome!
