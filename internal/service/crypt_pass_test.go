package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewCryptPass(t *testing.T) {
	cp := NewCryptPass()
	assert.NotNil(t, cp)
}

func TestHashPassword_Success(t *testing.T) {
	cp := NewCryptPass()
	password := "securePassword123"

	hashedPassword, err := cp.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	assert.NotEqual(t, password, hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err)
}

func TestComparePasswords_Success(t *testing.T) {
	cp := NewCryptPass()
	password := "securePassword123"

	hashedPassword, err := cp.HashPassword(password)
	assert.NoError(t, err)

	isMatch := cp.ComparePasswords(hashedPassword, password)
	assert.True(t, isMatch)
}

func TestComparePasswords_Failure(t *testing.T) {
	cp := NewCryptPass()
	password := "securePassword123"
	wrongPassword := "wrongPassword456"

	hashedPassword, err := cp.HashPassword(password)
	assert.NoError(t, err)

	isMatch := cp.ComparePasswords(hashedPassword, wrongPassword)
	assert.False(t, isMatch)
}
