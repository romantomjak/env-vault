package command

import (
	"fmt"
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

		file, err := os.CreateTemp(os.TempDir(), "env-vault-*")
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		filename := file.Name()
		defer os.Remove(filename)

		_, err = file.Write(plaintext)
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		if err = file.Close(); err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		if err = openFileInEditor(filename); err != nil {
			return fmt.Errorf("Error: %v", err)
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
	switch {
	case os.Getenv("ENV_VAULT_EDITOR") != "":
		return os.Getenv("ENV_VAULT_EDITOR")
	case os.Getenv("EDITOR") != "":
		return os.Getenv("EDITOR")
	default:
		return DefaultEditor
	}
}
