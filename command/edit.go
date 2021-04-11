package command

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/romantomjak/env-vault/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var editCmd = &cobra.Command{
	Use:   "edit [flags] filename",
	Short: "Edit encrypted file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if ENVCOMPOSE_PASSWORD_FILE was set or otherwise ask for password
		fmt.Fprintf(os.Stdin, "Password: ")
		bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println()

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

		tmpPlaintext, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		ciphertext2, err := vault.Encrypt(tmpPlaintext, bytepw)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(args[0], ciphertext2, 0700); err != nil {
			return err
		}

		return nil
	},
}
