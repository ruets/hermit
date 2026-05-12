package secrets

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ruets/hermit/internal/secrets/generators"
)

func Generate(name, secretType, secretsDir string) error {
	path := filepath.Join(secretsDir, name)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	gen, err := generators.New(secretType, name, secretsDir)
	if err != nil {
		return err
	}

	return gen.Generate(path)
}

// ExistsSecret checks if a secret exists in either encrypted or plaintext form
// based on the encrypted flag
func ExistsSecret(name, secretType string, secretsDir string, encrypted bool) bool {
	path := filepath.Join(secretsDir, name)

	if encrypted {
		agePath := path + ".age"

		// Check if main encrypted file exists
		info, err := os.Stat(agePath)
		return err == nil && info.Size() > 0
	}

	info, err := os.Stat(path)
	return err == nil && info.Size() > 0
}
