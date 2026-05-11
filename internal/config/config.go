package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Secret struct {
	Name     string   `yaml:"name"`
	Type     string   `yaml:"type"`
	Notes    string   `yaml:"notes"`
}

type Config struct {
	Secrets []Secret `yaml:"secrets"`
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
