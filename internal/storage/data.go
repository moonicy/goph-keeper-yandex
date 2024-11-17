package storage

import (
	"context"
	"errors"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	"gorm.io/gorm"
)

type DataRepository struct {
	db *gorm.DB
}

func NewDataRepository(db *gorm.DB) (*DataRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &DataRepository{db}, nil
}

func (dr *DataRepository) AddData(ctx context.Context, data entity.Data) error {
	db := dr.db.WithContext(ctx).Create(&data)
	if db.Error != nil {
		return db.Error
	}

	return db.Error
}
