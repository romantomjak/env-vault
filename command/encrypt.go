package command

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/vault"
)

var encryptCmd = &cobra.Command{
	Use:                   "encrypt [filename]",
	Short:                 "Permanently encrypt a plain text file",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := passwordPrompt("New password: ")
		if err != nil {
			return err
		}

		password2, err := passwordPrompt("Confirm new password: ")
		if err != nil {
			return err
		}

		if !bytes.Equal(password, password2) {
			return fmt.Errorf("passwords do not match")
		}

		plaintext, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}

		if err := vault.WriteFile(args[0], plaintext, password); err != nil {
			return err
		}

		return nil
	},
}
