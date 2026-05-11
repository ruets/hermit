package secrets

import (
	"fmt"
	"os"
	"path/filepath"

	agelib "filippo.io/age"
	"github.com/ruets/hermit/internal/config"
	"github.com/ruets/hermit/internal/secrets/age"
)

type Manager struct {
	cfg        *config.Config
	secretsDir string
	identity   *agelib.X25519Identity
}

func NewManager(cfg *config.Config, secretsDir string, identity *agelib.X25519Identity) *Manager {
	return &Manager{cfg: cfg, secretsDir: secretsDir, identity: identity}
}

type SecretStatus struct {
	Name     string
	Type     string
	Notes    string
	Exists   bool
}

func (m *Manager) Status() []SecretStatus {
	statuses := make([]SecretStatus, len(m.cfg.Secrets))
	for i, s := range m.cfg.Secrets {
		statuses[i] = SecretStatus{
			Name:     s.Name,
			Type:     s.Type,
			Notes:    s.Notes,
			Exists:   ExistsSecret(s.Name, s.Type, m.secretsDir, s.IsEncrypted()),
		}
	}
	return statuses
}

func (m *Manager) Generate() error {
	for _, s := range m.cfg.Secrets {
		if ExistsSecret(s.Name, s.Type, m.secretsDir, s.IsEncrypted()) {
			fmt.Printf("  ~ skipped   %s\n", s.Name)
			continue
		}
		if err := Generate(s.Name, s.Type, m.secretsDir); err != nil {
			return fmt.Errorf("failed to generate %s: %w", s.Name, err)
		}

		// Encrypt the generated secret(s) if needed
		if s.IsEncrypted() {
			path := filepath.Join(m.secretsDir, s.Name)
			if err := age.EncryptFile(m.identity, path); err != nil {
				return fmt.Errorf("failed to encrypt %s: %w", s.Name, err)
			}

			// For RSA keys, also encrypt the .pub file
			if s.Type == "rsa" {
				pubPath := path + ".pub"
				if err := age.EncryptFile(m.identity, pubPath); err != nil {
					return fmt.Errorf("failed to encrypt %s.pub: %w", s.Name, err)
				}
			}
		}

		fmt.Printf("  ✓ generated %s\n", s.Name)
	}
	return nil
}

func (m *Manager) Unwrap() error {
	secretsDir := ".secrets"

	for _, s := range m.cfg.Secrets {
		// Skip non-encrypted secrets
		if !s.IsEncrypted() {
			continue
		}

		path := filepath.Join(m.secretsDir, s.Name)
		agePath := path + ".age"

		// Check if encrypted file exists
		if _, err := os.Stat(agePath); err != nil {
			continue // Doesn't exist, skip
		}

		outPath := filepath.Join(secretsDir, s.Name)
		if err := age.DecryptFileTo(m.identity, agePath, outPath); err != nil {
			return fmt.Errorf("failed to decrypt %s: %w", s.Name, err)
		}
		fmt.Printf("  ✓ decrypted %s → .secrets/\n", s.Name)

		// For RSA keys, also decrypt the .pub file
		if s.Type == "rsa" {
			pubAgePath := path + ".pub.age"
			outPubPath := filepath.Join(secretsDir, s.Name+".pub")
			if err := age.DecryptFileTo(m.identity, pubAgePath, outPubPath); err != nil {
				return fmt.Errorf("failed to decrypt %s.pub: %w", s.Name, err)
			}
			fmt.Printf("  ✓ decrypted %s.pub → .secrets/\n", s.Name)
		}
	}
	return nil
}

func (m *Manager) Clean() error {
	secretsDir := ".secrets"

	for _, s := range m.cfg.Secrets {
		// Only clean plaintext files for encrypted secrets
		if !s.IsEncrypted() {
			continue
		}

		path := filepath.Join(secretsDir, s.Name)

		// Remove plaintext file from .secrets
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", s.Name, err)
		}

		// For RSA keys, also remove the .pub file
		if s.Type == "rsa" {
			pubPath := path + ".pub"
			if err := os.Remove(pubPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove %s.pub: %w", s.Name, err)
			}
		}
	}
	fmt.Println("  ✓ cleaned .secrets/ directory")
	return nil
}

func (m *Manager) Wrap() error {
	for _, s := range m.cfg.Secrets {
		// Only wrap secrets that should be encrypted
		if !s.IsEncrypted() {
			continue
		}

		path := filepath.Join(m.secretsDir, s.Name)
		agePath := path + ".age"

		// Skip if already encrypted
		if _, err := os.Stat(agePath); err == nil {
			continue
		}

		// Check if plaintext exists
		if _, err := os.Stat(path); err != nil {
			continue // Plaintext doesn't exist, skip
		}

		// Encrypt the plaintext file
		if err := age.EncryptFile(m.identity, path); err != nil {
			return fmt.Errorf("failed to encrypt %s: %w", s.Name, err)
		}

		// For RSA keys, also encrypt the .pub file
		if s.Type == "rsa" {
			pubPath := path + ".pub"
			if _, err := os.Stat(pubPath); err == nil {
				if err := age.EncryptFile(m.identity, pubPath); err != nil {
					return fmt.Errorf("failed to encrypt %s.pub: %w", s.Name, err)
				}
			}
		}

		fmt.Printf("  ✓ encrypted %s\n", s.Name)
	}
	return nil
}
