package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type CryptPass struct {
}

func NewCryptPass() *CryptPass {
	return &CryptPass{}
}

// HashPassword хэширует пароль с использованием bcrypt
func (CryptPass) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// ComparePasswords сравнивает хэш пароля с введённым паролем
func (CryptPass) ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
