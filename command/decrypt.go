package command

import (
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/vault"
)

var decryptCmd = &cobra.Command{
	Use:                   "decrypt [filename]",
	Short:                 "Permanently decrypt an encrypted file",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := passwordFromEnvOrPrompt("Password: ")
		if err != nil {
			return err
		}

		plaintext, err := vault.ReadFile(args[0], password)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(args[0], plaintext, 0700); err != nil {
			return err
		}

		return nil
	},
}
