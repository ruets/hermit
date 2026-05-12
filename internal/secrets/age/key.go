package age

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
)

func GenerateKey(keyPath string) (*age.X25519Identity, error) {
	if err := os.MkdirAll(filepath.Dir(keyPath), 0700); err != nil {
		return nil, fmt.Errorf("failed to create key directory: %w", err)
	}

	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, fmt.Errorf("failed to generate age key: %w", err)
	}

	if err := os.WriteFile(keyPath, []byte(identity.String()+"\n"), 0600); err != nil {
		return nil, fmt.Errorf("failed to write key: %w", err)
	}

	return identity, nil
}

func LoadKey(keyPath string) (*age.X25519Identity, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("key not found at %s: %w", keyPath, err)
	}

	identities, err := age.ParseIdentities(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}

	return identities[0].(*age.X25519Identity), nil
}
