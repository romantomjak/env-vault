package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

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
