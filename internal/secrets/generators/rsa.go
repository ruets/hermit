package generators

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type RSAGenerator struct{}

func (g *RSAGenerator) Generate(path string) error {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	})

	if err := writeSecret(path, string(privPEM)); err != nil {
		return err
	}

	return writeSecret(path+".pub", string(pubPEM))
}
