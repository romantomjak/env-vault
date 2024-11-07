package command

import (
	"fmt"
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

		modifiedplaintext, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		if err := vault.WriteFile(args[0], modifiedplaintext, password); err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		return nil
	},
}
