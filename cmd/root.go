package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ruets/hermit/internal/config"
	"github.com/ruets/hermit/internal/secrets"
	"github.com/ruets/hermit/internal/secrets/age"
	"github.com/spf13/cobra"
)

var (
	configPath string
	secretsDir string
	keyDir     string
)

var rootCmd = &cobra.Command{
	Use:           "hermit",
	Short:         "A tool to manage your secrets from a config file",
	SilenceUsage:  true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	homeDir, _ := os.UserHomeDir()
	defaultKeyDir := filepath.Join(homeDir, ".config", "hermit")

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "secrets.yaml", "path to secrets.yaml")
	rootCmd.PersistentFlags().StringVarP(&secretsDir, "secrets-dir", "d", "secrets", "path to secrets directory")
	rootCmd.PersistentFlags().StringVarP(&keyDir, "key-dir", "k", defaultKeyDir, "path to age key directory")
}

func newManager() (*secrets.Manager, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Determine key directory: config > flag > default
	resolvedKeyDir := keyDir
	if cfg.KeyDir != "" {
		resolvedKeyDir = cfg.KeyDir
	}

	// Expand ~ in path
	if resolvedKeyDir[0] == '~' {
		homeDir, _ := os.UserHomeDir()
		resolvedKeyDir = filepath.Join(homeDir, resolvedKeyDir[1:])
	}

	// Load or generate age key
	identity, err := age.LoadKey(resolvedKeyDir)
	if err != nil {
		// Key doesn't exist, generate a new one
		identity, err = age.GenerateKey(resolvedKeyDir)
		if err != nil {
			return nil, fmt.Errorf("failed to generate age key: %w", err)
		}
		keyPath := age.KeyPath(resolvedKeyDir)
		fmt.Printf("✓ generated age key at %s\n", keyPath)
		fmt.Println("⚠ back up this file — losing it means losing access to all secrets")
	}

	return secrets.NewManager(cfg, secretsDir, identity), nil
}
