package generators

import (
	"fmt"
	"os"
)

type Generator interface {
	Generate(path string) error
}

func New(secretType, name, secretsDir string) (Generator, error) {
	switch secretType {
	case "random_hex":
		return &RandomHexGenerator{}, nil
	case "rsa":
		return &RSAGenerator{}, nil
	case "manual":
		return &ManualGenerator{name: name}, nil
	default:
		return nil, fmt.Errorf("unknown secret type: %s", secretType)
	}
}

func writeSecret(path, value string) error {
	if err := os.WriteFile(path, []byte(value), 0600); err != nil {
		return fmt.Errorf("failed to write secret: %w", err)
	}
	return nil
}
