package command

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/vault"
)

var editCmd = &cobra.Command{
	Use:                   "edit [filename]",
	Short:                 "Edit encrypted file",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			return err
		}

		password, err := passwordFromEnvOrPrompt("Password: ")
		if err != nil {
			return err
		}

		plaintext, err := vault.ReadFile(args[0], password)
		if err != nil {
			return err
		}

		file, err := ioutil.TempFile(os.TempDir(), "env-vault-*")
		if err != nil {
			return err
		}

		filename := file.Name()
		defer os.Remove(filename)

		_, err = file.Write(plaintext)
		if err != nil {
			return err
		}

		if err = file.Close(); err != nil {
			return err
		}

		if err = openFileInEditor(filename); err != nil {
			return err
		}

		modifiedplaintext, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		if err := vault.WriteFile(args[0], modifiedplaintext, password); err != nil {
			return err
		}

		return nil
	},
}
