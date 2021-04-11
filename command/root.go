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

var rootCmd = &cobra.Command{
	Use:   "env-vault [vault] [command]",
	Short: "Launch a subprocess with environment variables from an encrypted file",
	Long: `env-vault provides a convenient way to launch a subprocess
with environment variables populated from an encrypted file.`,
	Args:                  cobra.MinimumNArgs(2),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		executable, err := exec.LookPath(args[1])
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdin, "Password: ")
		bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println()

		vaultbytes, err := vault.ReadFile(args[0], bytepw)
		if err != nil {
			return err
		}

		// TODO: add flag for stripping current env vars
		newEnv := make([]string, 0)
		newEnv = append(newEnv, os.Environ()...)
		newEnv = append(newEnv, strings.Split(string(vaultbytes), "\n")...)

		c := exec.Command(executable, args[2:]...)
		c.Env = newEnv
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		return c.Run()
	},
}

func Execute() error {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(editCmd)

	return rootCmd.Execute()
}
