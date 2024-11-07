package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/vault"
)

var decryptCmd = &cobra.Command{
	Use:                   "decrypt [filename]",
	Short:                 "Permanently decrypt an encrypted file",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			return fmt.Errorf("Error: %v", err)
		}

		password, err := passwordFromEnvOrPrompt("Password: ")
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		plaintext, err := vault.ReadFile(args[0], password)
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		if err := os.WriteFile(args[0], plaintext, 0700); err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		return nil
	},
}
