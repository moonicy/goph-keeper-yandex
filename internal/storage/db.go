package storage

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB создает новое соединение с базой данных.
func NewDB(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("no Database specified")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db, nil
}
