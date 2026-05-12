package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	configPath string
	keyPath    string
)

var rootCmd = &cobra.Command{
	Use:          "hermit",
	Short:        "A tool to manage your secrets from a config file",
	SilenceUsage: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	homeDir, _ := os.UserHomeDir()
	defaultKeyPath := filepath.Join(homeDir, ".config", "hermit", "hermit.key")

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "secrets.yaml", "path to secrets.yaml")
	rootCmd.PersistentFlags().StringVarP(&keyPath, "key-path", "k", defaultKeyPath, "path to age key file")

	rootCmd.AddGroup(&cobra.Group{
		ID:    "management",
		Title: "Management Commands:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "encryption",
		Title: "Encryption Commands:",
	})

}
