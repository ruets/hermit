# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Hermit** is a CLI tool for managing secrets from a YAML configuration file. It handles encryption/decryption using the `age` cryptography library, generation of various secret types (random hex, RSA keys), and provides a workflow for safely working with secrets locally.

## Common Commands

### Build & Run
```bash
go build -o hermit .
./hermit --help
```

### Development
```bash
go vet ./...          # Lint check
go fmt ./...          # Format code
go mod tidy           # Clean dependencies
```

## Architecture

### High-Level Workflow

The hermit workflow has 5 steps:

1. **Generate** (`hermit generate`)
   - Reads `secrets.yaml` configuration
   - Creates secrets in `secretsDir` (default: `./secrets/`)
   - Encrypts with `.age` suffix if `encrypted: true` (default)
   - For RSA keys, generates private key (public key derived on-the-fly during unwrap)

2. **Unwrap** (`hermit unwrap`)
   - Decrypts encrypted secrets from `secretsDir` to `.secrets/` (plaintext)
   - Generates RSA public keys on-the-fly from private keys
   - Used before user modifications

3. **User Edits** (manual)
   - Secrets are temporarily in plaintext in `.secrets/`
   - User can modify them

4. **Wrap** (`hermit wrap`)
   - Scans `.secrets/` and compares each file with stored version
   - Detects modifications (changed content)
   - Asks user confirmation for modified secrets
   - Re-encrypts and saves to `secretsDir`
   - Removes plaintext from `.secrets/` (or deletes `.secrets/` if empty)

5. **Clean** (`hermit clean`)
   - Scans both `.secrets/` and `secretsDir` for orphaned files
   - Files that don't exist in `secrets.yaml` are orphaned
   - Asks user confirmation before deletion

### Package Structure

```
internal/
├── config/          # YAML configuration parsing
│   └── config.go    # Config struct, Secret types
├── secrets/
│   ├── manager.go   # Main logic (Generate, Unwrap, Clean, Wrap, Status)
│   ├── helpers.go   # Utility functions (confirm, writeWithBackup)
│   ├── generate.go  # Secret generation orchestration
│   ├── age/
│   │   ├── crypto.go    # Encrypt/Decrypt bytes (age.Encrypt, age.Decrypt)
│   │   ├── file.go      # File I/O (EncryptFile, DecryptFile, DecryptFileTo)
│   │   └── key.go       # Key loading/generation
│   └── generators/      # Secret type generators
│       ├── interface.go # Generator interface
│       ├── random_hex.go
│       ├── rsa.go
│       └── manual.go

cmd/
├── root.go          # Root cobra command + flags (--config, --key-path)
├── init.go          # `hermit init` command
├── generate.go      # `hermit generate` command
├── clean.go         # `hermit clean` command
├── status.go        # `hermit status` command
├── wrap.go          # `hermit wrap` command
└── unwrap.go        # `hermit unwrap` command
```

### Key Design Principles

- **Manager-centric**: `internal/secrets/Manager` is the main orchestrator that:
  - Holds the age identity and secrets directory
  - Implements all workflow steps (Generate, Unwrap, Clean, Wrap)
  - Calls helpers and generators as needed

- **Separation of concerns**:
  - `age/` package: Pure cryptography operations (encrypt/decrypt bytes and files)
  - `helpers.go`: Business logic helpers (save with backup, compare encrypted content)
  - `generators/`: Pluggable secret generators
  - `Manager`: Coordinates between config, helpers, and generators

- **RSA handling**: For RSA secrets:
  - Only the private key is stored and encrypted
  - Public key is mathematically derived from the private key (using `crypto/x509`)
  - Public key is generated on-the-fly during `unwrap()` to `.secrets/` for user access
  - `.pub` files are never stored permanently, avoiding duplication and storage overhead

## Configuration

`secrets.yaml` format:
```yaml
secrets:
  - name: authelia/jwt_secret    # Secret name (directory structure supported)
    type: random_hex             # Type: random_hex, rsa, manual, or custom
    notes: authelia              # Optional comma-separated tags
    encrypted: true              # Optional, defaults to true
```

Available types:
- `random_hex`: Generate random hexadecimal string
- `rsa`: Generate RSA private/public key pair
- `manual`: User provides the secret value interactively
- Custom generators can be added in `internal/secrets/generators/`

## CLI Flags

Available on all commands (set in `cmd/root.go`):
```bash
hermit [command] \
  --config secrets.yaml \
  --key-path ~/.config/hermit/hermit.key
```

**Path Standardization**: Paths are deterministic and calculated from the config location:
- `secretsDir` = `{configDir}/secrets/` (encrypted secrets)
- `.secrets/` = `{configDir}/.secrets/` (temporary plaintext for user editing)
- The `--secrets-dir` flag is intentionally removed to ensure paths are predictable

## Important Implementation Notes

### Wrap & Clean Logic

Both `Wrap()` and `Clean()` use `filepath.WalkDir` to scan directories:
- **Wrap**: Scans `.secrets/` for plaintext files, compares with stored versions
  - For encrypted secrets: decrypts stored `.age` file in-memory for comparison
  - For RSA secrets: only compares private keys (public keys are derived on-the-fly)
  - Detects modifications and asks user confirmation before re-encrypting
- **Clean**:
  - Phase A: Scans `.secrets/` for orphaned files (not in `secrets.yaml`)
  - Phase B: Scans `secretsDir` for orphaned `.age` files (not in `secrets.yaml`)
  - RSA public keys in `.secrets/` are never orphaned since they're generated on-the-fly and always regenerated

### Encryption Details

- Uses `filippo.io/age` (modern encryption standard)
- `age.Encrypt(identity, plaintext)` returns encrypted bytes
- `age.Decrypt(identity, ciphertext)` returns plaintext bytes
- Files are written/read with `os.WriteFile` / `os.ReadFile`
- Encrypted files get `.age` suffix appended
- RSA public keys (`.pub`) are generated from private keys using `crypto/x509.MarshalPKCS1PublicKey()`
- Public keys are PEM-encoded and never stored permanently (generated on-the-fly as needed)

### Backup Strategy

When saving modified secrets:
- Existing file is renamed to `.bak` (backup)
- New content is written
- No automatic cleanup of `.bak` files

### User Confirmation

The `confirm()` helper asks for `[y/N]` (default is No):
```go
func confirm(prompt string) bool {
    // Reads from stdin, case-insensitive
    // Returns true only for "y" or "yes"
}
```

## Testing

Test data in `test/`:
- `secrets.yaml`: Example configuration with auth services (authelia, lldap, etc)
- `secrets/`: Directory structure with encrypted secrets
- `.secrets/`: Decrypted plaintext for testing
