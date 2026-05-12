package cmd

import (
	"github.com/ruets/hermit/internal/secrets"
	"github.com/spf13/cobra"
)

var unwrapCmd = &cobra.Command{
	Use:     "unwrap",
	GroupID: "encryption",
	Short:   "Decrypt all secrets to plaintext",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := secrets.NewManager(configPath, keyPath)
		if err != nil {
			return err
		}
		return m.Unwrap()
	},
}

func init() {
	rootCmd.AddCommand(unwrapCmd)
}
