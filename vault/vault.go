package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"os"
)

type Vault struct {
	file     *os.File
	password string
}

func Open(filename, password string) (*Vault, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		return nil, err
	}

	// TODO: check file header

	v := &Vault{
		file:     f,
		password: password,
	}
	return v, nil
}

func (v *Vault) Read() ([]byte, error) {
	cipherText, err := ioutil.ReadAll(v.file)
	if err != nil {
		return nil, err
	}
	return Decrypt(cipherText, []byte(v.password))
}

func (v *Vault) Write(b []byte) (n int, err error) {
	cipherText, err := Encrypt(b, []byte(v.password))
	if err != nil {
		return 0, err
	}
	return v.file.Write(cipherText)
}

func (v *Vault) Close() error {
	return v.file.Close()
}

// cipher key should be 32 bit long, so lets generate one by hashing password
func cipherKey(password []byte) []byte {
	key := sha256.Sum256(password)
	return key[:]
}

func Encrypt(bytes, password []byte) ([]byte, error) {
	block, err := aes.NewCipher(cipherKey(password))
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, bytes, nil)

	return ciphertext, nil
}

func Decrypt(bytes, password []byte) ([]byte, error) {
	block, err := aes.NewCipher(cipherKey(password))
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := bytes[:nonceSize], bytes[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
