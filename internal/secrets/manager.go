package secrets

import (
	"fmt"

	"github.com/ruets/hermit/internal/config"
)

type Manager struct {
	cfg        *config.Config
	secretsDir string
}

func NewManager(cfg *config.Config, secretsDir string) *Manager {
	return &Manager{cfg: cfg, secretsDir: secretsDir}
}

type SecretStatus struct {
	Name     string
	Type     string
	Services []string
	Exists   bool
}

