package command

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/vault"
)

var createCmd = &cobra.Command{
	Use:                   "create [filename]",
	Short:                 "Create new encrypted file",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := passwordPrompt("Password: ")
		if err != nil {
			return err
		}

		password2, err := passwordPrompt("Confirm new vault password: ")
		if err != nil {
			return err
		}

		if !bytes.Equal(password, password2) {
			return fmt.Errorf("passwords do not match")
		}

		file, err := ioutil.TempFile(os.TempDir(), "env-vault-*")
		if err != nil {
			return err
		}

		filename := file.Name()
		defer os.Remove(filename)

		if err = file.Close(); err != nil {
			return err
		}

		if err = openFileInEditor(filename); err != nil {
			return err
		}

		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		if err := vault.WriteFile(args[0], bytes, password); err != nil {
			return err
		}

		return nil
	},
}
