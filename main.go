package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const DefaultEditor = "vim"

func getPreferredEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return DefaultEditor
	}
	return editor
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

func decrypt(enc, password []byte) (string, error) {
	// cipher key should be 32 bit long, so lets generate one by hashing password
	key := sha256.Sum256(password)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func main() {
	createCmd := &cobra.Command{
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

			// cipher key should be 32 bit long, so lets generate one by hashing password
			key := sha256.Sum256(bytepw)

			block, err := aes.NewCipher(key[:])
			if err != nil {
				return err
			}

			aesGCM, err := cipher.NewGCM(block)
			if err != nil {
				return err
			}

			nonce := make([]byte, aesGCM.NonceSize())
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				return err
			}

			ciphertext := aesGCM.Seal(nonce, nonce, bytes, nil)

			if err := ioutil.WriteFile(args[0], ciphertext, 0700); err != nil {
				return err
			}

			return nil
		},
	}

	viewCmd := &cobra.Command{
		Use:   "view [flags] filename",
		Short: "View encrypted file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// check if ENVCOMPOSE_PASSWORD_FILE was set or otherwise ask for password
			fmt.Fprintf(os.Stdin, "Password: ")
			bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}

			// cipher key should be 32 bit long, so lets generate one by hashing password
			key := sha256.Sum256(bytepw)

			block, err := aes.NewCipher(key[:])
			if err != nil {
				return err
			}

			aesGCM, err := cipher.NewGCM(block)
			if err != nil {
				return err
			}

			enc, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			nonceSize := aesGCM.NonceSize()

			nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

			plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
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

	rootCmd := &cobra.Command{
		Use:   "env-vault [OPTIONS] [COMMAND]",
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

			filename := "secrets.yml"
			enc, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}

			envs, err := decrypt(enc, bytepw)
			if err != nil {
				return err
			}

			c := exec.Command(executable)
			c.Env = strings.Split(envs, "\n")
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			return c.Run()
		},
	}
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(viewCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// set ENVCOMPOSE_PASSWORD_FILE to automatically search for password in that file
	// create to create new encrypted file
	// encrypt to encryp existing file
	// view encrypted file
	// edit encrypted file
	// rekey to change password of encrypted file
	// decrypt to permanently decrypt
}
