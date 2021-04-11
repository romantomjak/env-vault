package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/romantomjak/env-vault/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var VaultFile string

var rootCmd = &cobra.Command{
	Use:   "env-vault [flags] [command]",
	Short: "Launch a subprocess with environment variables from an encrypted file",
	Long: `env-vault provides a convenient way to launch a subprocess
with environment variables populated from an encrypted file.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executable, err := exec.LookPath(args[0])
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdin, "Password: ")
		bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println()

		v, err := vault.Open(VaultFile, string(bytepw))
		if err != nil {
			return err
		}

		plaintext, err := v.Read()
		if err != nil {
			return err
		}

		// TODO: add flag for stripping current env vars
		newEnv := make([]string, 0)
		newEnv = append(newEnv, os.Environ()...)
		newEnv = append(newEnv, strings.Split(string(plaintext), "\n")...)

		c := exec.Command(executable, args[1:]...)
		c.Env = newEnv
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		return c.Run()
	},
}

func Execute() error {
	rootCmd.Flags().StringVar(&VaultFile, "vault", "secrets.env", "filename to read secrets from")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(editCmd)

	return rootCmd.Execute()
}
