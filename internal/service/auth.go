package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
)

type AuthService struct {
	userRepository UserRepository
	cryptPass      ICryptPass
	tokenGenerator ITokenGenerator
}

func NewAuthService(userRepository UserRepository, cryptPass ICryptPass, tokenGenerator ITokenGenerator) (*AuthService, error) {
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

func (as *AuthService) Login(ctx context.Context, login string, password string) (string, string, error) {
	if login == "" {
		return "", "", errors.New("login is empty")
	}
	if password == "" {
		return "", "", errors.New("password is empty")
	}
	user, err := as.userRepository.Get(ctx, login)
	if err != nil {
		return "", "", fmt.Errorf("user %s is not exist", login)
	}
	if !as.cryptPass.ComparePasswords(user.Password, password) {
		return "", "", errors.New("wrong password")
	}
	token, err := as.tokenGenerator.GenerateToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("generating token: %w", err)
	}
	return token, user.Salt, nil
}

func (as *AuthService) Register(ctx context.Context, login string, password string) (entity.User, error) {
	cryptPass, err := as.cryptPass.HashPassword(password)
	if err != nil {
		return entity.User{}, err
	}
	return as.userRepository.Create(ctx, entity.User{
		Login:    login,
		Password: cryptPass,
		Salt:     uuid.New().String(),
	})
}
