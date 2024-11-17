package storage

import (
	"errors"
	"github.com/moonicy/goph-keeper-yandex/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB создает новое соединение с базой данных.
func NewDB(cfg config.ServerConfig) (*gorm.DB, error) {
	if cfg.Database == "" {
		return nil, errors.New("no Database specified")
	}
	db, err := gorm.Open(postgres.Open(cfg.Database), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db, nil
}
