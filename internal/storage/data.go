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

func (dr *DataRepository) RemoveData(ctx context.Context, id uint64) error {
	db := dr.db.WithContext(ctx).Delete(&entity.Data{}, id)
	if db.Error != nil {
		return db.Error
	}
	return db.Error
}

func (dr *DataRepository) GetData(ctx context.Context, userID uint64) ([]entity.Data, error) {
	var data []entity.Data
	db := dr.db.WithContext(ctx).Where("user_id = ?", userID).Find(&data)
	if db.Error != nil {
		return []entity.Data{}, db.Error
	}
	return data, nil
}

func (dr *DataRepository) UpdateData(ctx context.Context, data entity.Data) error {
	db := dr.db.WithContext(ctx).Save(&data)
	if db.Error != nil {
		return db.Error
	}
	return db.Error
}
