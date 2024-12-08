package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCrypt(t *testing.T) {
	crypt := NewCrypt()
	assert.NotNil(t, crypt)
}

func TestCrypt_Init(t *testing.T) {
	crypt := NewCrypt()

	password := "securePassword123"
	salt := "randomSalt456"

	crypt.Init(password, salt)

	assert.NotNil(t, crypt)
}

func TestCrypt_EncryptDecrypt_Success(t *testing.T) {
	crypt := NewCrypt()

	password := "securePassword123"
	salt := "randomSalt456"
	crypt.Init(password, salt)

	plaintext := []byte("Hello, this is a test message!")

	encryptedData, err := crypt.Encrypt(plaintext)
	assert.NoError(t, err)
	assert.NotNil(t, encryptedData)

	decryptedData, err := crypt.Decrypt(encryptedData)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decryptedData)
}

func TestCrypt_Decrypt_ErrorInvalidData(t *testing.T) {
	crypt := NewCrypt()

	password := "securePassword123"
	salt := "randomSalt456"
	crypt.Init(password, salt)

	invalidData := []byte("invalid encrypted data")

	_, err := crypt.Decrypt(invalidData)
	assert.Error(t, err)
}

func TestCrypt_Clean(t *testing.T) {
	crypt := NewCrypt()

	password := "securePassword123"
	salt := "randomSalt456"
	crypt.Init(password, salt)

	crypt.Clean()

	assert.Nil(t, crypt.encryptionKey)
}
