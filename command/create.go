package command

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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
