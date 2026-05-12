package age

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
)

const keyFile = "hermit.key"

func KeyPath(configDir string) string {
	return filepath.Join(configDir, keyFile)
}

func GenerateKey(configDir string) (*age.X25519Identity, error) {
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config dir: %w", err)
	}

	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, fmt.Errorf("failed to generate age key: %w", err)
	}

	path := KeyPath(configDir)
	if err := os.WriteFile(path, []byte(identity.String()+"\n"), 0600); err != nil {
		return nil, fmt.Errorf("failed to write key: %w", err)
	}

	return identity, nil
}

func LoadKey(configDir string) (*age.X25519Identity, error) {
	path := KeyPath(configDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("key not found at %s: %w", path, err)
	}

	identities, err := age.ParseIdentities(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}

	return identities[0].(*age.X25519Identity), nil
}
