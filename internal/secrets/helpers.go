package secrets

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/ruets/hermit/internal/secrets/age"
)

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("  ? %s [y/N] ", prompt)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

// writeWithBackup writes content to path with backup, handling both plaintext and encrypted
func (m *Manager) writeWithBackup(path string, content []byte, encrypt bool) error {
	// Backup existing file(s)
	if _, err := os.Stat(path); err == nil {
		os.Rename(path, path+".bak")
	}
	if _, err := os.Stat(path + ".age"); err == nil {
		os.Rename(path+".age", path+".age.bak")
	}

	// Write new content
	if encrypt {
		ciphertext, err := age.Encrypt(m.identity, content)
		if err != nil {
			return fmt.Errorf("failed to encrypt: %w", err)
		}
		return os.WriteFile(path+".age", ciphertext, 0400)
	}
	return os.WriteFile(path, content, 0644)
}

// generateRSAPublicKey derives and writes the public key from private key bytes.
func generateRSAPublicKey(secretName string, privKeyBytes []byte, publicKeyPath string) error {
	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		return fmt.Errorf("failed to decode PEM for %s", secretName)
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key %s: %w", secretName, err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privKey.PublicKey),
	})
	if err := os.WriteFile(publicKeyPath, pubPEM, 0644); err != nil {
		return fmt.Errorf("failed to write public key for %s: %w", secretName, err)
	}

	return nil
}
