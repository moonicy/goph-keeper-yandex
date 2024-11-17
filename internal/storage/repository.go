package storage

import (
	"errors"
	"github.com/moonicy/goph-keeper-yandex/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// BaseRepository представляет базу данных с возможностью повторных попыток при ошибках соединения.
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository создает новое соединение с базой данных.
func NewBaseRepository(cfg config.ServerConfig) (*BaseRepository, error) {
	if cfg.Database == "" {
		return nil, errors.New("no Database specified")
	}
	db, err := gorm.Open(postgres.Open(cfg.Database), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return &BaseRepository{db: db}, nil
}
