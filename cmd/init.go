package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	GroupID: "management",
	Short:   "Initialize hermit in a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if secrets.yaml already exists
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("⚠ Project already initialized: %s exists\n", configPath)
			return nil
		}

		// Create secrets.yaml
		exampleConfig := `# Optional: specify the age encryption key path
# If not specified, defaults to ~/.config/hermit/hermit.key
# key_path: ./hermit.key

secrets:
  # Random hex secrets (tokens, API keys, etc.)
  - name: random_hex_example
    type: random_hex
    notes: example random hex secret

  # RSA key pairs (private key stored, public key derived)
  - name: rsa_example
    type: rsa
    notes: example RSA key pair

  # Unencrypted secret (optional, default is encrypted)
  - name: manual_example
    type: manual
    encrypted: false
    notes: example manual secret
`

		if err := os.WriteFile(configPath, []byte(exampleConfig), 0o600); err != nil {
			return fmt.Errorf("failed to create %s: %w", configPath, err)
		}
		fmt.Printf("✓ created %s\n", configPath)

		// Update .gitignore
		gitignorePath := ".gitignore"
		gitignoreEntry := "# Hermit secrets (NEVER commit encryption keys!)\n.secrets/\n*.key\nhermit.key\n"

		// Read existing .gitignore if it exists
		var content []byte
		if _, err := os.Stat(gitignorePath); err == nil {
			content, err = os.ReadFile(gitignorePath)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", gitignorePath, err)
			}
		}

		// Check if .secrets/ is already in .gitignore
		contentStr := string(content)
		if strings.Contains(contentStr, ".secrets/") {
			fmt.Printf("✓ .secrets/ already ignored in .gitignore\n")
		} else {
			// Append to .gitignore
			if len(content) > 0 && !strings.HasSuffix(contentStr, "\n") {
				content = append(content, '\n')
			}
			content = append(content, []byte(gitignoreEntry)...)

			if err := os.WriteFile(gitignorePath, content, 0o644); err != nil {
				return fmt.Errorf("failed to write %s: %w", gitignorePath, err)
			}
			fmt.Printf("✓ updated .gitignore (added .secrets/)\n")
		}

		// Print next steps
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  1. Edit %s to configure your secrets\n", configPath)
		fmt.Println("  2. hermit generate         # Generate the secrets")
		fmt.Println("  3. git add secrets/ .gitignore")
		fmt.Println("  4. git commit -m 'init: add hermit secrets'")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
