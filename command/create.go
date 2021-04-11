package command

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/romantomjak/env-vault/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var createCmd = &cobra.Command{
	Use:   "create [flags] filename",
	Short: "Create new encrypted file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if ENVCOMPOSE_PASSWORD_FILE was set or otherwise ask for password
		fmt.Fprintf(os.Stdin, "New vault password: ")
		bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
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

		v, err := vault.Open(args[0], string(bytepw))
		if err != nil {
			return err
		}

		if _, err := v.Write(bytes); err != nil {
			return err
		}

		if err = v.Close(); err != nil {
			return err
		}

		return nil
	},
}
