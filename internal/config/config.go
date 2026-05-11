package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Secret struct {
	Name     	string   `yaml:"name"`
	Type     	string   `yaml:"type"`
	Notes    	string   `yaml:"notes"`
	Encrypted *bool    `yaml:"encrypted"` // nil = default to true
}

type Config struct {
	KeyDir  string   `yaml:"key_dir"`  // optional, defaults to ~/.config/hermit
	Secrets []Secret `yaml:"secrets"`
}

func (s *Secret) IsEncrypted() bool {
	if s.Encrypted == nil {
		return true // default to encrypted
	}
	return *s.Encrypted
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
