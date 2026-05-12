package cmd

import (
	"github.com/ruets/hermit/internal/secrets"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:     "clean",
	GroupID: "management",
	Short:   "Remove orphaned secrets with confirmation",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := secrets.NewManager(configPath, keyPath)
		if err != nil {
			return err
		}
		return m.Clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
