package generators

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type RandomHexGenerator struct{}

func (g *RandomHexGenerator) Generate(path string) error {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return writeSecret(path, hex.EncodeToString(b))
}
