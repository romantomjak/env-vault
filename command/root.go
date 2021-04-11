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

func passwordPrompt(prompt string, stdin, stdout *os.File) ([]byte, error) {
	// check if ENVCOMPOSE_PASSWORD_FILE was set or otherwise ask for password

	fmt.Fprint(stdout, prompt)

	password, err := term.ReadPassword(int(stdin.Fd()))
	if err != nil {
		return nil, err
	}

	// ReadPassword skips newline
	fmt.Fprintln(stdout)

	return password, nil
}

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

		password, err := passwordPrompt("Password: ", os.Stdin, os.Stdout)
		if err != nil {
			return err
		}

		vaultbytes, err := vault.ReadFile(args[0], password)
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
