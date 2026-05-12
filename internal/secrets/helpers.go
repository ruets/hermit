package secrets

import (
	"bufio"
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
		return os.WriteFile(path+".age", ciphertext, 0600)
	}
	return os.WriteFile(path, content, 0600)
}
