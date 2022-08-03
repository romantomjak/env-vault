package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/romantomjak/env-vault/vault"
)

func passwordFromEnvOrPrompt(prompt string) ([]byte, error) {
	envpassword := os.Getenv("ENV_VAULT_PASSWORD")
	if len(envpassword) > 0 {
		return []byte(envpassword), nil
	}
	return passwordPrompt(prompt)
}

func passwordPrompt(prompt string) ([]byte, error) {
	fmt.Fprint(os.Stdout, prompt)

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}

	// add newline skipped by term.ReadPassword
	fmt.Fprintln(os.Stdout)

	return password, nil
}

// pristineEnv is a command line flag that controls whether
// the sub-process will inherit environment variables from
// the current process. Defaults to false
var pristineEnv bool

var rootCmd = &cobra.Command{
	Use:   "env-vault [options] vault program",
	Short: "Launch a subprocess with environment variables from an encrypted file",
	Long: `env-vault provides a convenient way to launch a subprocess
with environment variables populated from an encrypted file.`,
	Args:                  cobra.MinimumNArgs(2),
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	SilenceErrors:         true,
	RunE: func(cmd *cobra.Command, args []string) error {
		executable, err := exec.LookPath(args[1])
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		password, err := passwordFromEnvOrPrompt("Password: ")
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		vaultbytes, err := vault.ReadFile(args[0], password)
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		newEnv := make([]string, 0)

		// by default sub-process will inherit current process environment
		// variables unless specified otherwise
		if !pristineEnv {
			newEnv = append(newEnv, os.Environ()...)
		}

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
	rootCmd.Flags().BoolVar(&pristineEnv, "pristine", false, "Only use values retrieved from vaults, do not inherit the existing environment variables")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(rekeyCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}
