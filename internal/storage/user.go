package storage

import (
	"context"
	"errors"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) (*UserRepository, error) {
	if db == nil {
		return nil, errors.New("NewUserRepository: db is nil")
	}
	return &UserRepository{db: db}, nil
}

func (ur *UserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	db := ur.db.WithContext(ctx).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Create(&user)
	if db.Error != nil {
		return entity.User{}, db.Error
	}

	return user, db.Error
}

func (ur *UserRepository) Get(ctx context.Context, login string) (entity.User, error) {
	var user entity.User
	db := ur.db.WithContext(ctx).First(&user, "login = ?", login)
	return user, db.Error
}
