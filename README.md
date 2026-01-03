# Environment Vault (ev)

**ev** is a secure, lightweight CLI for managing environment variables. It works similarly to `ansible-vault` but is designed for speed, portability, and easy integration with shell scripts.

Built in Go, `ev` compiles to a single binary with no dependencies.

## Features

- 🔒 **AES-256-GCM Encryption**: Your secrets are safe.
- 🔑 **PBKDF2 Key Derivation**: Protected against brute-force attacks.
- 📝 **Base64 Encoding**: Vault files are Git-friendly and copy-pasteable.
- 🐚 **Shell Native**: First-class support for exporting variables to your shell.
- 🚀 **Zero Dependencies**: Just a single binary.

## Comparison: `ev` vs. `ansible-vault`

While `ansible-vault` is a powerful general-purpose tool, `ev` is purpose-built for environment variable management in shell scripts and CI/CD pipelines.

| Feature | ev 🚀 | ansible-vault 🐢 |
| :--- | :--- | :--- |
| **Startup Speed** | **~20ms** (Instant) | **~500ms+** (Python startup) |
| **Dependencies** | **None** (Single Binary) | Python, Ansible, Pip modules |
| **Shell Integration** | Native (`ev export`, `ev get`) | Difficult (requires parsing output) |
| **File Format** | Base64 (Copy-paste friendly) | Raw Binary / Hex |
| **Editing** | Auto-cleanup of temp files | Auto-cleanup of temp files |
| **Encryption** | AES-256-GCM (Authenticated) | AES-256 (CBC/CTR) |

**Why use ev?**
- **For Scripts:** running `eval $(ev export)` is instant and loads your entire environment safely. Doing this with Ansible requires complex piping and is significantly slower.
- **For CI/CD:** You don't need to install a heavy Python environment just to unlock secrets.

## Installation

### Homebrew (macOS & Linux)
The recommended way to install. This use a custom tap to keep `ev` up-to-date.

```bash
brew install hugovntr/tap/ev
```

To update later:
```bash
brew upgrade ev
```

### Go Install
If you are a Go developer, you can install directly from the source:
```bash
go install github.com/hugovntr/ev@latest
```

### Manual Binary
1. Download the latest binary for your OS/Arch from the [Release page](http://github.com/hugovntr/ev/releases)
2. Decompress the archive
3. Move it to your path:
```bash
tar -xzf ev_*.tar.gz
sudo mv ev /usr/local/bin
```

### From Source
Requires Go 1.22+

```bash
git clone https://github.com/hugovntr/ev.git
cd ev
make install
```

## Verification
Verify the installation by checking the version
```bash
ev --version
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

4. **Load Environment**
Load all vault variables into your current shell.

**Bash / Zsh**
```bash
eval $(ev env)
```

**Fish**
```fish
ev env fish | source
```

**PowerShell**
```powershell
ev env pwsh | Invoke-Expression
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
