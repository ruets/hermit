package cmd

import (
	"github.com/ruets/hermit/internal/secrets"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	GroupID: "management",
	Short:   "Generate all missing secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := secrets.NewManager(configPath, keyPath)
		if err != nil {
			return err
		}
		return m.Generate()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
