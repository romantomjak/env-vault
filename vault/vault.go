package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

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
