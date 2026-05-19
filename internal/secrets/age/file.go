package age

import (
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
)

func EncryptFile(identity *age.X25519Identity, path string) error {
	plaintext, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ciphertext, err := Encrypt(identity, plaintext)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path+".age", ciphertext, 0400); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	return os.Remove(path)
}

func DecryptFile(identity *age.X25519Identity, path string) error {
	ciphertext, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	plaintext, err := Decrypt(identity, ciphertext)
	if err != nil {
		return err
	}

	outPath := path[:len(path)-4] // supprime .age
	return os.WriteFile(outPath, plaintext, 0644)
}

func DecryptFileTo(identity *age.X25519Identity, path, outPath string) error {
	ciphertext, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	plaintext, err := Decrypt(identity, ciphertext)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(outPath, plaintext, 0644)
}
