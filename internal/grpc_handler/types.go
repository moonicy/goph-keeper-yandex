package grpc_handler

import (
	"context"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
)

//go:generate mockery --output ./../../mocks --filename auth_servise_gen.go --outpkg mocks --name AuthService --with-expecter
type AuthService interface {
	Login(ctx context.Context, login string, password string) (string, string, error)
	Register(ctx context.Context, login string, password string) (entity.User, error)
}

//go:generate mockery --output ./../../mocks --filename data_repository_gen.go --outpkg mocks --name DataRepository --with-expecter
type DataRepository interface {
	AddData(ctx context.Context, data entity.Data) error
	RemoveData(ctx context.Context, id uint64) error
	GetData(ctx context.Context, userID uint64) ([]entity.Data, error)
	UpdateData(ctx context.Context, data entity.Data) error
}
