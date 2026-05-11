package cmd

import (
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate all missing secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := newManager()
		if err != nil {
			return err
		}
		return m.Generate()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
