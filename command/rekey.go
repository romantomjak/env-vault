package command

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/vault"
)

var rekeyCmd = &cobra.Command{
	Use:                   "rekey [filename]",
	Short:                 "Change password of an encrypted file",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			return fmt.Errorf("Error: %v", err)
		}

		password, err := passwordFromEnvOrPrompt("Current password: ")
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		newpassword, err := passwordPrompt("New password: ")
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		newpassword2, err := passwordPrompt("Confirm new password: ")
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		if !bytes.Equal(newpassword, newpassword2) {
			return fmt.Errorf("Error: passwords do not match")
		}

		plaintext, err := vault.ReadFile(args[0], password)
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		if err := vault.WriteFile(args[0], plaintext, newpassword); err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		return nil
	},
}
