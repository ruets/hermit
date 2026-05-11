package age

import (
	"bytes"
	"fmt"
	"io"
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

func Encrypt(identity *age.X25519Identity, plaintext []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, identity.Recipient())
	if err != nil {
		return nil, fmt.Errorf("failed to create age writer: %w", err)
	}

	if _, err := w.Write(plaintext); err != nil {
		return nil, fmt.Errorf("failed to encrypt: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize encryption: %w", err)
	}

	return buf.Bytes(), nil
}

func Decrypt(identity *age.X25519Identity, ciphertext []byte) ([]byte, error) {
	r, err := age.Decrypt(bytes.NewReader(ciphertext), identity)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	plaintext, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read decrypted data: %w", err)
	}

	return plaintext, nil
}

func EncryptFile(identity *age.X25519Identity, path string) error {
	plaintext, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ciphertext, err := Encrypt(identity, plaintext)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path+".age", ciphertext, 0600); err != nil {
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
	return os.WriteFile(outPath, plaintext, 0600)
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

	return os.WriteFile(outPath, plaintext, 0600)
}
