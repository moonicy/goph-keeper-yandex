package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

type Crypt struct {
	encryptionKey []byte
}

func NewCrypt() *Crypt {
	return &Crypt{}
}

func (c *Crypt) Init(password string, salt string) {
	c.encryptionKey = pbkdf2.Key([]byte(password), []byte(salt), 10000, 32, sha256.New)
}

// Encrypt Метод для шифрования данных
func (c *Crypt) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	encryptedData := append(nonce, ciphertext...)
	return encryptedData, nil
}

// Decrypt Метод для дешифровки данных
func (c *Crypt) Decrypt(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("шифротекст слишком короткий")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (c *Crypt) Clean() {
	c.encryptionKey = nil
}
