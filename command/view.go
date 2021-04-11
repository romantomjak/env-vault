package command

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/vault"
)

const DefaultEditor = "vim"

var viewCmd = &cobra.Command{
	Use:                   "view [filename]",
	Short:                 "View encrypted file",
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

		return nil
	},
}

func openFileInEditor(filename string) error {
	executable, err := exec.LookPath(getPreferredEditor())
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getPreferredEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return DefaultEditor
	}
	return editor
}
