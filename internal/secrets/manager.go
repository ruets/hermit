package secrets

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	agelib "filippo.io/age"
	"github.com/ruets/hermit/internal/config"
	"github.com/ruets/hermit/internal/secrets/age"
)

type Manager struct {
	cfg        *config.Config
	secretsDir string
	identity   *agelib.X25519Identity
}

func NewManager(configPath, keyPath string) (*Manager, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Calculate secretsDir from configPath
	configDir := filepath.Dir(configPath)
	secretsDir := filepath.Join(configDir, "secrets")

	// Resolve keyPath
	resolvedKeyPath := keyPath
	if len(resolvedKeyPath) > 0 && resolvedKeyPath[0] == '~' {
		homeDir, _ := os.UserHomeDir()
		resolvedKeyPath = filepath.Join(homeDir, resolvedKeyPath[1:])
	}

	identity, err := age.LoadKey(resolvedKeyPath)
	if err != nil {
		identity, err = age.GenerateKey(resolvedKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to generate age key: %w", err)
		}
		fmt.Printf("✓ generated age key at %s\n", resolvedKeyPath)
		fmt.Println("⚠ back up this file — losing it means losing access to all secrets")
	}

	return &Manager{cfg: cfg, secretsDir: secretsDir, identity: identity}, nil
}

type SecretStatus struct {
	Name   string
	Type   string
	Notes  string
	Exists bool
}

func (m *Manager) Status() []SecretStatus {
	statuses := make([]SecretStatus, len(m.cfg.Secrets))
	for i, s := range m.cfg.Secrets {
		statuses[i] = SecretStatus{
			Name:   s.Name,
			Type:   s.Type,
			Notes:  s.Notes,
			Exists: ExistsSecret(s.Name, s.Type, m.secretsDir, s.IsEncrypted()),
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

		// Encrypt the generated secret if needed
		if s.IsEncrypted() {
			path := filepath.Join(m.secretsDir, s.Name)
			if err := age.EncryptFile(m.identity, path); err != nil {
				return fmt.Errorf("failed to encrypt %s: %w", s.Name, err)
			}
		}

		fmt.Printf("  ✓ generated %s\n", s.Name)
	}
	return nil
}

func (m *Manager) Clean() error {
	tempSecretsDir := filepath.Join(filepath.Dir(m.secretsDir), ".secrets")
	scanDirs := []string{tempSecretsDir, m.secretsDir}
	orphans := make(map[string]bool)

	// Collect all orphaned secrets
	for _, scanDir := range scanDirs {
		if err := filepath.WalkDir(scanDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				if os.IsNotExist(err) {
					return filepath.SkipAll
				}
				return err
			}

			if d.IsDir() {
				return nil
			}

			rel, _ := filepath.Rel(scanDir, path)

			// Ignore .pub files (handled with their parent)
			if strings.HasSuffix(rel, ".pub") {
				return nil
			}

			// Extract secret name
			var secretName string
			if strings.HasSuffix(rel, ".age") {
				secretName = strings.TrimSuffix(rel, ".age")
			} else {
				secretName = rel
			}

			// Check if this secret exists in config
			found := false
			for _, s := range m.cfg.Secrets {
				if s.Name == secretName {
					found = true
					break
				}
			}

			if !found {
				orphans[secretName] = true
			}

			return nil
		}); err != nil {
			return err
		}
	}

	// Ask for confirmation and delete
	for orphan := range orphans {
		if confirm(fmt.Sprintf("delete orphaned secret %s?", orphan)) {
			// Delete from both directories
			for _, scanDir := range scanDirs {
				// Try plaintext version
				mainFile := filepath.Join(scanDir, orphan)
				os.Remove(mainFile)
				// Try encrypted version
				mainFile = filepath.Join(scanDir, orphan+".age")
				os.Remove(mainFile)

				// Try pub versions
				pubFile := filepath.Join(scanDir, orphan+".pub")
				os.Remove(pubFile)
				pubFile = filepath.Join(scanDir, orphan+".pub.age")
				os.Remove(pubFile)
			}
			fmt.Printf("  ✓ deleted %s\n", orphan)
		}
	}

	fmt.Println("  ✓ cleaned secrets")
	return nil
}

func (m *Manager) Unwrap() error {
	tempSecretsDir := filepath.Join(filepath.Dir(m.secretsDir), ".secrets")

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

		outPath := filepath.Join(tempSecretsDir, s.Name)
		if err := age.DecryptFileTo(m.identity, agePath, outPath); err != nil {
			return fmt.Errorf("failed to decrypt %s: %w", s.Name, err)
		}
		fmt.Printf("  ✓ decrypted %s → .secrets/\n", s.Name)

		// For RSA keys, generate the public key as well
		if s.Type == "rsa" {
			privPEM, err := os.ReadFile(outPath)
			if err != nil {
				return fmt.Errorf("failed to read private key: %w", err)
			}

			block, _ := pem.Decode(privPEM)
			if block == nil {
				return fmt.Errorf("failed to decode PEM for %s", s.Name)
			}

			privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse private key %s: %w", s.Name, err)
			}

			pubPEM := pem.EncodeToMemory(&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(&privKey.PublicKey),
			})
			outPubPath := filepath.Join(tempSecretsDir, s.Name+".pub")
			if err := os.WriteFile(outPubPath, pubPEM, 0600); err != nil {
				return fmt.Errorf("failed to decrypt %s.pub: %w", s.Name, err)
			}
			fmt.Printf("  ✓ generated %s.pub → .secrets/\n", s.Name)
		}
	}
	return nil
}

func (m *Manager) Wrap() error {
	tempSecretsDir := filepath.Join(filepath.Dir(m.secretsDir), ".secrets")
	type change struct {
		name       string
		secret     config.Secret
		content    []byte
		isModified bool
	}
	changes := make(map[string]*change)

	// Collect all changes
	if err := filepath.WalkDir(tempSecretsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return filepath.SkipAll
			}
			return err
		}

		if d.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(tempSecretsDir, path)

		// Ignore .pub files (handled with their parent)
		if strings.HasSuffix(rel, ".pub") {
			return nil
		}

		// Find secret in config
		var secret *config.Secret
		for i, s := range m.cfg.Secrets {
			if s.Name == rel {
				secret = &m.cfg.Secrets[i]
				break
			}
		}

		if secret == nil {
			return nil // Skip, let clean handle it
		}

		// Read new content
		newContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", rel, err)
		}

		// Get stored content
		storedPath := filepath.Join(m.secretsDir, secret.Name)
		var storedContent []byte

		if secret.IsEncrypted() {
			agePath := storedPath + ".age"
			if _, err := os.Stat(agePath); err == nil {
				ciphertext, err := os.ReadFile(agePath)
				if err != nil {
					return fmt.Errorf("failed to read %s: %w", secret.Name, err)
				}
				decrypted, err := age.Decrypt(m.identity, ciphertext)
				if err != nil {
					return fmt.Errorf("failed to decrypt %s: %w", secret.Name, err)
				}
				storedContent = decrypted
			}
		} else {
			if _, err := os.Stat(storedPath); err == nil {
				storedContent, err = os.ReadFile(storedPath)
				if err != nil {
					return fmt.Errorf("failed to read %s: %w", secret.Name, err)
				}
			}
		}

		// Check if modified
		isModified := !bytes.Equal(newContent, storedContent)

		changes[secret.Name] = &change{
			name:       rel,
			secret:     *secret,
			content:    newContent,
			isModified: isModified,
		}

		return nil
	}); err != nil {
		return err
	}

	// Process changes
	for _, ch := range changes {
		if ch.isModified {
			// Ask for modified files
			if confirm(fmt.Sprintf("%s modified — save?", ch.name)) {
				// Save main file
				storedPath := filepath.Join(m.secretsDir, ch.secret.Name)
				if err := m.writeWithBackup(storedPath, ch.content, ch.secret.IsEncrypted()); err != nil {
					return fmt.Errorf("failed to save %s: %w", ch.secret.Name, err)
				}

				fmt.Printf("  ✓ saved %s\n", ch.name)
			} else {
				fmt.Printf("  ~ discarded %s (not saved)\n", ch.name)
			}
		}

		// Remove from .secrets
		mainPath := filepath.Join(tempSecretsDir, ch.name)
		if err := os.Remove(mainPath); err == nil || os.IsNotExist(err) {
			if !ch.isModified {
				fmt.Printf("  ~ removed %s (unchanged)\n", ch.name)
			}
		}
		if ch.secret.Type == "rsa" {
			pubPath := filepath.Join(tempSecretsDir, ch.name+".pub")
			os.Remove(pubPath)
		}
	}

	// Try to remove .secrets/ directory if empty
	if err := os.Remove(tempSecretsDir); err == nil {
		fmt.Println("  ✓ removed .secrets/ (empty)")
	}

	return nil
}
