package command

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/vault"
)

var encryptCmd = &cobra.Command{
	Use:                   "encrypt [filename]",
	Short:                 "Permanently encrypt a plain text file",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			return fmt.Errorf("Error: %v", err)
		}

		password, err := passwordFromEnvOrPrompt("New password: ")
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		password2, err := passwordFromEnvOrPrompt("Confirm new password: ")
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		if !bytes.Equal(password, password2) {
			return fmt.Errorf("Error: passwords do not match")
		}

		plaintext, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		if err := vault.WriteFile(args[0], plaintext, password); err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		return nil
	},
}
