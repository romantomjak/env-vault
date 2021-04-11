package command

import (
	"io/ioutil"
	"os"

	"github.com/romantomjak/env-vault/vault"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [flags] filename",
	Short: "Edit encrypted file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bytepw, err := passwordPrompt("Password: ", os.Stdin, os.Stdout)
		if err != nil {
			return err
		}

		plaintext, err := vault.ReadFile(args[0], bytepw)
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

		if err := vault.WriteFile(args[0], modifiedplaintext, bytepw); err != nil {
			return err
		}

		return nil
	},
}
