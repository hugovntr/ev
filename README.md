# EV (Environment Vault)

**ev** is a secure, lightweight CLI for managing environment variables. It works similarly to `ansible-vault` but is designed for speed, portability, and easy integration with shell scripts.

Built in Go, `ev` compiles to a single binary with no dependencies.

## Features

- 🔒 **AES-256-GCM Encryption**: Your secrets are safe.
- 🔑 **PBKDF2 Key Derivation**: Protected against brute-force attacks.
- 📝 **Base64 Encoding**: Vault files are Git-friendly and copy-pasteable.
- 🐚 **Shell Native**: First-class support for exporting variables to your shell.
- 🚀 **Zero Dependencies**: Just a single binary.

## Installation

### From Source
Requires Go 1.22+

```bash
git clone https://github.com/hugovntr/ev.git
cd ev
make install
```

## Usage
1. **Initialize a Vault**
Create a new encrypted vault. You will be prompted to create a password.
```bash
ev create
```

*Default locations*:
- Vault: `~/.config/ev/vault`
- Credentials: `~/.config/ev/credentials`

2. **Edit Variables**
Opens your default `$EDITOR` (vim, nano, code). The file is decrypted for editing and strictly re-encrypted upon save.

```bash
ev edit
```

*Format* (YAML):
```YAML
API_KEY: "super-secret"
DB_HOST: "localhost"
DEBUG: true
```

3. **Get a Value**
Retrieve a single decrypted value. Useful for piping into other commands.

```bash
# Example: Using the token in a curl request
curl -H "Authorization: Bearer $(ev get API_KEY)" https://api.example.com
```

4. **Export All Variables**
Load all vault variables into your current shell session.

```bash
eval $(ev export)
```

### Configuration Flags
You can override the default paths using flags or by editing the source defaults.

- `-f`, `--vault-file`: Path to the encrypted vault file.
- `-p`, `--password-file`: Path to the file containing the vault password.

## Security
**ev** uses a robust security model:
1. **Key Derivation**: Passwords are salted (16 bytes) and hashed using PBKDF2 (SHA-256, 100,000 iterations).
2. **Encryption**: Data is encrypted using AES-256-GCM.
3. **Integrity**: GCM ensures that if the ciphertext is tampered with, decryption will fail immediately.
4. **Ephemeral Files**: Plaintext data exists on disk only momentarily during ev edit and is securely deleted immediately after the editor closes.

## License

```text
MIT License

Copyright (c) Hugo Ventura

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
