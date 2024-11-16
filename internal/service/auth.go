package service

import (
	"context"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	"github.com/moonicy/goph-keeper-yandex/internal/storage"
)

type AuthService struct {
	userRepository storage.UserRepository
}

func NewAuthService(userRepository storage.UserRepository) *AuthService {
	return &AuthService{userRepository}
}

func (as *AuthService) Login(ctx context.Context, login string, password string) error {
	return nil
}

func (as *AuthService) Logout() error {
	return nil
}

func (as *AuthService) Register(ctx context.Context, login string, password string) error {
	return as.userRepository.Create(ctx, entity.User{
		Login:    login,
		Password: password,
	})
}
