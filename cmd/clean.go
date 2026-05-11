package cmd

import (
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove plaintext secrets, keep encrypted .age files",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := newManager()
		if err != nil {
			return err
		}
		return m.Clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
