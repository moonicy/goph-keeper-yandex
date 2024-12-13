package service

import (
	"context"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
)

//go:generate mockery --output ./../../mocks --filename user_repo_gen.go --outpkg mocks --name UserRepository --with-expecter
type UserRepository interface {
	Get(ctx context.Context, login string) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
}

//go:generate mockery --output ./../../mocks --filename crypt_pass_gen.go --outpkg mocks --name ICryptPass --with-expecter
type ICryptPass interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword, password string) bool
}

//go:generate mockery --output ./../../mocks --filename token_gen_gen.go --outpkg mocks --name ITokenGenerator --with-expecter
type ITokenGenerator interface {
	GenerateToken(userID uint64) (string, error)
}
