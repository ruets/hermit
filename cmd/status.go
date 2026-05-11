package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	GroupID: "management",
	Short: "List all secrets and their status",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := newManager()
		if err != nil {
			return err
		}
		for _, s := range m.Status() {
			status := "✓"
			if !s.Exists {
				status = "✗"
			}
			fmt.Printf("  %s  %-40s %s\n", status, s.Name, s.Notes)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
