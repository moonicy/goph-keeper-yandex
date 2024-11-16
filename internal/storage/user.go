package storage

import (
	"context"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
)

type UserRepository struct {
	repo BaseRepository
}

func NewUserRepository(repo BaseRepository) *UserRepository {
	return &UserRepository{repo: repo}
}

func (ur *UserRepository) Create(ctx context.Context, user entity.User) error {
	return ur.repo.db.WithContext(ctx).Create(&user).Error
}
