package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

func Encrypt(bytes, password []byte) ([]byte, error) {
	// cipher key should be 32 bit long, so lets generate one by hashing password
	key := sha256.Sum256(password)

	block, err := aes.NewCipher(key[:])
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

func Decrypt(enc, password []byte) (string, error) {
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
