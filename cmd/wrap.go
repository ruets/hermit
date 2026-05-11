package cmd

import (
	"github.com/spf13/cobra"
)

var wrapCmd = &cobra.Command{
	Use:     "wrap",
	GroupID: "encryption",
	Short:   "Encrypt plaintext secrets marked as encrypted in config",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := newManager()
		if err != nil {
			return err
		}
		return m.Wrap()
	},
}

func init() {
	rootCmd.AddCommand(wrapCmd)
}
