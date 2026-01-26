package tradingaccount

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type AESCrypto struct {
	gcm cipher.AEAD
}

func NewAESCrypto(key []byte) (*AESCrypto, error) {
	if len(key) != 32 {
		return nil, errors.New("AES key must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AESCrypto{gcm: gcm}, nil
}

func (a *AESCrypto) Encrypt(plain string) (string, error) {
	nonce := make([]byte, a.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := a.gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (a *AESCrypto) Decrypt(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	nonceSize := a.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plain, err := a.gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
