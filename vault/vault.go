package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"io/ioutil"
)

func ReadFile(filename string, password []byte) ([]byte, error) {
	ciphertext, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return decrypt(ciphertext, password)
}

func WriteFile(filename string, data, password []byte) error {
	ciphertext, err := encrypt(data, password)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, ciphertext, 0700)
}

// cipher key should be 32 bit long, so lets generate one by hashing password
func cipherKey(password []byte) []byte {
	key := sha256.Sum256(password)
	return key[:]
}

func encrypt(bytes, password []byte) ([]byte, error) {
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

func decrypt(bytes, password []byte) ([]byte, error) {
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
