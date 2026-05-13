# Hermit

**Hermit** is a minimalist secret manager designed to let infrastructure admins version their secrets securely. Secrets are generated and stored encrypted alongside their configuration, forming a single versioned source of truth.

## Why This Project?

### The Problem

When managing a homelab with many services (Authelia, Immich, LLDAP, Vikunja, Outline, etc.), you face a classic dilemma:

- **Option 1**: Store secrets in plaintext → security risk
- **Option 2**: Generate secrets manually and store them separately → config/secrets drift
- **Option 3**: Use a complex secret manager like Vault → overkill for a homelab

### The Solution: Hermit

Hermit offers a simple but effective approach:

1. **Declare** secrets in `secrets.yaml` alongside your config
2. **Generate** secrets automatically (RSA keys, random values, etc.)
3. **Encrypt** secrets with version them in Git
4. **Modify** secrets securely by decrypting them locally
5. **Re-encrypt** before committing

**Key advantage**: Your configuration and its secrets live together, versioned, encrypted, without depending on external infrastructure.

### Use Case: Homelab with Multiple Services

Consider a typical homelab:

```
services:
  - Authelia        (centralized auth)
  - LLDAP           (user directory)
  - Forgejo         (self-hosted git)
```

Each service needs secrets:

- JWT keys for authentication
- Database passwords
- Encryption keys
- OAuth OIDC secrets

With Hermit, you organize everything in `secrets.yaml`:

```yaml
secrets:
  - name: authelia/jwt_secret
    type: random_hex

  - name: authelia/oidc
    type: rsa

  - name: authelia/oidc_clients/vikunja
    type: random_hex

  - name: lldap/jwt_secret
    type: random_hex
```

Then generate, encrypt, and version everything:

```bash
hermit generate          # Creates secrets encrypted
git commit -m "chore: generate secrets"
```

Later, to modify a secret:

```bash
hermit unwrap       # Decrypts to .secrets/ (Git-ignored)
vim .secrets/authelia/jwt_secret
hermit wrap         # Re-encrypts modifications
```

## Features

- **Simple**: Single YAML file to declare and configure secrets
- **Portable**: Encryption with age (modern standard)
- **Flexible**: Support for multiple secret types
  - `random_hex`: Random hexadecimal values
  - `rsa`: RSA key pairs (private + public)
  - `manual`: Interactive user input
  - Extensible with custom generators
- **Secure**: Secrets encrypted at rest, safe modification workflows
- **Organized**: Directory hierarchy supported (`authelia/jwt_secret`, `oidc_clients/vikunja`, etc.)
- **Versionable**: Encrypted secrets ready for Git, just backup your age key securely
- **No external state**: No Vault server, no network dependencies

## Installation

### Requirements

- Go 1.21+

### From Source

```bash
git clone https://github.com/ruets/hermit
cd hermit
go build -o hermit .
./hermit --help
```

## Configuration

Create a `secrets.yaml` file in your config directory or use `hermit init` to generate a default one:

```yaml
key_path: ~/.config/hermit/hermit.key   # Optional, path to age key

secrets:
  - name: authelia/jwt_secret
    type: random_hex

  - name: authelia/oidc
    type: rsa

  - name: authelia/oidc_clients/vikunja
    type: random_hex

  - name: lldap/jwt_secret
    type: random_hex
```

### Configuration Options

You can configure hermit via:

**YAML config file** (`secrets.yaml`):
- `key_path`: Path to age encryption key (optional, overrides CLI flag if set)

**CLI flags**:
```bash
hermit [command] \
  --config path/to/secrets.yaml \
  --key-path ~/.config/hermit/hermit.key
```

- `--config`: Path to `secrets.yaml` (default: `./secrets.yaml`)
- `--key-path`: Path to age key (default: `~/.config/hermit/hermit.key`)

**Priority**: If `key_path` is set in `secrets.yaml`, it takes priority over the `--key-path` flag.

Secrets are always stored alongside the config file:

- `{config_dir}/secrets/`: Encrypted secrets
- `{config_dir}/.secrets/`: Temporary plaintext (Git-ignored)

## Commands

| Command           | Description                    |
| ----------------- | ------------------------------ |
| `hermit init`     | Initialize Hermit config files |
| `hermit generate` | Generate missing secrets       |
| `hermit status`   | Show secret status             |
| `hermit unwrap`   | Decrypt secrets (plaintext)    |
| `hermit wrap`     | Re-encrypt modifications       |
| `hermit clean`    | Remove orphaned secrets        |

### `init`

```bash
hermit init
```

Initializes Hermit by creating a default `secrets.yaml` and adding secrets and key to gitignore.

### `generate`

```bash
hermit generate
```

Creates missing secrets in `{config_dir}/secrets/` (encrypted by default)

### `status`

```bash
hermit status
```

Shows list of secrets and their status

### `unwrap`

```bash
hermit unwrap
```

Decrypts secrets to `{config_dir}/.secrets/` (plaintext files)

For RSA keys, the public key is automatically generated from the private key

### `wrap`

```bash
hermit wrap
```

Re-encrypts secrets, save modified secrets and removes plaintext versions

### `clean`

```bash
hermit clean
```

Removes secrets no longer in `secrets.yaml`

## Supported Secret Types

### `random_hex`

Generates a random hexadecimal string.

```yaml
- name: shared/api_token
  type: random_hex
```

Generates: `a3f7e2c9d1b4e8f6c2a9e7d4b1f8c5a3`

### `rsa`

Generates a 2048-bit RSA key pair (private + public).

```yaml
- name: authelia/oidc
  type: rsa
```

Generates two files:

- `secrets/authelia/oidc.age`: Encrypted private key
- `.secrets/authelia/oidc` + `.pub`: Plaintext during decryption

The public key is **generated on-the-fly** from the private key, never stored.

### `manual`

Prompts user for interactive input.

```yaml
- name: shared/admin_password
  type: manual
```

## Security

### Encryption

- **Algorithm**: Age (filippo.io/age)
- **Key derivation**: X25519
- **Format**: `.age` files containing encrypted secrets

### Secure Workflow

1. Secrets encrypted at rest (`secrets/`)
2. Decryption happens locally only (`.secrets/`)
3. Plaintext editing on your machine
4. Re-encryption before commit
5. `.secrets/` ignored by Git

### Key Backup

```bash
# First time, hermit generates your age key
✓ generated age key at ~/.config/hermit/hermit.key
⚠ back up this file — losing it means losing access to all secrets
```

**Important**: Backup your age key (`~/.config/hermit/hermit.key`). Losing it means losing access to all encrypted secrets.

## Further Documentation

For detailed information about the architecture and development, see [CLAUDE.md](./CLAUDE.md).

## AI Usage Note

This project was primarily developed by hand, with AI assistance on specific parts. The AI contributed to:

- Initial code structure and some implementation details
- Documentation writing
- Testing and validation

The code has been reviewed and tested. Feel free to audit it for your use case, as with any open-source project.

## Contributing

Contributions welcome! Open an issue or pull request.

---

**Hermit**: Manage your secrets simply, keeping them close to your config, encrypted and versioned.
