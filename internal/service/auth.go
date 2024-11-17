package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	"github.com/moonicy/goph-keeper-yandex/internal/storage"
)

type AuthService struct {
	userRepository *storage.UserRepository
	cryptPass      *CryptPass
	tokenGenerator *TokenGenerator
}

func NewAuthService(userRepository *storage.UserRepository, cryptPass *CryptPass, tokenGenerator *TokenGenerator) (*AuthService, error) {
	if userRepository == nil {
		return nil, errors.New("userRepository is nil")
	}
	if cryptPass == nil {
		return nil, errors.New("cryptPass is nil")
	}
	return &AuthService{
		userRepository: userRepository,
		cryptPass:      cryptPass,
		tokenGenerator: tokenGenerator,
	}, nil
}

func (as *AuthService) Login(ctx context.Context, login string, password string) (string, error) {
	if login == "" {
		return "", errors.New("login is empty")
	}
	if password == "" {
		return "", errors.New("password is empty")
	}
	user, err := as.userRepository.Get(ctx, login)
	if err != nil {
		return "", fmt.Errorf("user %s is not exist", login)
	}
	if !as.cryptPass.ComparePasswords(user.Password, password) {
		return "", errors.New("wrong password")
	}
	token, err := as.tokenGenerator.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}
	return token, nil
}

func (as *AuthService) Logout() error {
	return nil
}

func (as *AuthService) Register(ctx context.Context, login string, password string) (uint64, error) {
	cryptPass, err := as.cryptPass.HashPassword(password)
	if err != nil {
		return 0, err
	}
	return as.userRepository.Create(ctx, entity.User{
		Login:    login,
		Password: cryptPass,
	})
}
