package cmd

import (
	"fmt"
	"os"

	"github.com/ruets/hermit/internal/config"
	"github.com/ruets/hermit/internal/secrets"
	"github.com/spf13/cobra"
)

var (
	configPath string
	secretsDir string
)

var rootCmd = &cobra.Command{
	Use:   "hermit",
	Short: "A tool to manage your secrets from a config file",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "secrets.yaml", "path to secrets.yaml")
	rootCmd.PersistentFlags().StringVarP(&secretsDir, "secrets-dir", "d", "secrets", "path to secrets directory")
}

func newManager() (*secrets.Manager, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return secrets.NewManager(cfg, secretsDir), nil
}
