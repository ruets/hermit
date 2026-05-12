package cmd

import (
	"github.com/ruets/hermit/internal/secrets"
	"github.com/spf13/cobra"
)

var wrapCmd = &cobra.Command{
	Use:     "wrap",
	GroupID: "encryption",
	Short:   "Clean .secrets/ and re-encrypt modified secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := secrets.NewManager(configPath, keyPath)
		if err != nil {
			return err
		}
		return m.Wrap()
	},
}

func init() {
	rootCmd.AddCommand(wrapCmd)
}
