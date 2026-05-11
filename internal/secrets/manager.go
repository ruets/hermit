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
	Notes    string
	Exists   bool
}

func (m *Manager) Status() []SecretStatus {
	statuses := make([]SecretStatus, len(m.cfg.Secrets))
	for i, s := range m.cfg.Secrets {
		statuses[i] = SecretStatus{
			Name:     s.Name,
			Type:     s.Type,
			Notes:    s.Notes,
			Exists:   Exists(s.Name, m.secretsDir),
		}
	}
	return statuses
}

func (m *Manager) Generate() error {
	for _, s := range m.cfg.Secrets {
		if Exists(s.Name, m.secretsDir) {
			fmt.Printf("  ~ skipped   %s\n", s.Name)
			continue
		}
		if err := Generate(s.Name, s.Type, m.secretsDir); err != nil {
			return fmt.Errorf("failed to generate %s: %w", s.Name, err)
		}
		fmt.Printf("  ✓ generated %s\n", s.Name)
	}
	return nil
}
