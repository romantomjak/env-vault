package command

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/vault"
)

var rekeyCmd = &cobra.Command{
	Use:                   "rekey [filename]",
	Short:                 "Re-key an encrypted file",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := passwordFromEnvOrPrompt("Current password: ")
		if err != nil {
			return err
		}

		newpassword, err := passwordPrompt("New password: ")
		if err != nil {
			return err
		}

		newpassword2, err := passwordPrompt("Confirm new password: ")
		if err != nil {
			return err
		}

		if !bytes.Equal(newpassword, newpassword2) {
			return fmt.Errorf("new passwords do not match")
		}

		plaintext, err := vault.ReadFile(args[0], password)
		if err != nil {
			return err
		}

		if err := vault.WriteFile(args[0], plaintext, newpassword); err != nil {
			return err
		}

		return nil
	},
}
