package cmd

import (
	"github.com/spf13/cobra"
)

var unwrapCmd = &cobra.Command{
	Use:   "unwrap",
	Short: "Decrypt all secrets to plaintext",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := newManager()
		if err != nil {
			return err
		}
		return m.Unwrap()
	},
}

func init() {
	rootCmd.AddCommand(unwrapCmd)
}
