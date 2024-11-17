package storage

import (
	"context"
	"errors"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	repo *BaseRepository
}

func NewUserRepository(repo *BaseRepository) (*UserRepository, error) {
	if repo == nil {
		return nil, errors.New("NewUserRepository: repository is nil")
	}
	return &UserRepository{repo: repo}, nil
}

func (ur *UserRepository) Create(ctx context.Context, user entity.User) (uint, error) {
	db := ur.repo.db.WithContext(ctx).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Create(&user)
	if db.Error != nil {
		return 0, db.Error
	}

	return user.ID, db.Error
}

func (ur *UserRepository) Get(ctx context.Context, login string) (entity.User, error) {
	var user entity.User
	db := ur.repo.db.WithContext(ctx).First(&user, "login = ?", login)
	return user, db.Error
}
