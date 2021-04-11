package vault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func ReadFile(filename string, password []byte) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	firstNewLine := bytes.IndexByte(data, '\n')
	if firstNewLine < 0 {
		return nil, fmt.Errorf("vault header not found")
	}

	header := data[0:firstNewLine]
	if err := checkHeader(header); err != nil {
		return nil, err
	}

	body := data[firstNewLine+1:]

	// decode base64 encoded bytes
	ciphertext := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
	n, err := base64.StdEncoding.Decode(ciphertext, body)
	if err != nil {
		return nil, err
	}
	ciphertext = ciphertext[0:n]

	return decrypt(ciphertext, password)
}

func WriteFile(filename string, data, password []byte) error {
	ciphertext, err := encrypt(data, password)
	if err != nil {
		return err
	}

	header := "env-vault;1.0;AES256\n"

	// encode bytes as base64
	body := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(body, ciphertext)

	buf := bytes.NewBufferString(header)
	buf.Write(body)
	return ioutil.WriteFile(filename, buf.Bytes(), 0700)
}

func checkHeader(data []byte) error {
	header := string(data)
	lines := strings.SplitN(header, ";", 3)
	if lines[0] != "env-vault" {
		return fmt.Errorf("unknown format ID. was the file encrypted with env-vault?")
	}
	if lines[1] != "1.0" {
		return fmt.Errorf("incompatible file version. only 1.0 is supported for now")
	}
	if lines[2] != "AES256" {
		return fmt.Errorf("unsupported cipher algorithm. only AES256 is supported for now")
	}
	return nil
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
