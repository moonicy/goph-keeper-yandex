package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/google/uuid"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	"github.com/moonicy/goph-keeper-yandex/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthService_Success(t *testing.T) {
	userRepoMock := &mocks.UserRepository{}
	cryptPassMock := &mocks.ICryptPass{}
	tokenGenMock := &mocks.ITokenGenerator{}

	authService, err := NewAuthService(userRepoMock, cryptPassMock, tokenGenMock)

	assert.NoError(t, err)
	assert.NotNil(t, authService)
}

func TestNewAuthService_Error(t *testing.T) {
	cryptPassMock := &mocks.ICryptPass{}
	tokenGenMock := &mocks.ITokenGenerator{}

	authService, err := NewAuthService(nil, cryptPassMock, tokenGenMock)
	assert.Error(t, err)
	assert.Nil(t, authService)
	assert.Equal(t, "userRepository is nil", err.Error())

	authService, err = NewAuthService(&mocks.UserRepository{}, nil, tokenGenMock)
	assert.Error(t, err)
	assert.Nil(t, authService)
	assert.Equal(t, "cryptPass is nil", err.Error())
}

func TestAuthService_Login_Success(t *testing.T) {
	ctx := context.Background()
	login := "testuser"
	password := "testpassword"
	userID := uint64(1)
	salt := uuid.New().String()
	token := "mockToken"

	userRepoMock := &mocks.UserRepository{}
	cryptPassMock := &mocks.ICryptPass{}
	tokenGenMock := &mocks.ITokenGenerator{}

	userRepoMock.EXPECT().Get(ctx, login).Return(entity.User{
		ID:       userID,
		Login:    login,
		Password: "hashedpassword",
		Salt:     salt,
	}, nil)

	cryptPassMock.EXPECT().ComparePasswords("hashedpassword", password).Return(true)
	tokenGenMock.EXPECT().GenerateToken(userID).Return(token, nil)

	authService, err := NewAuthService(userRepoMock, cryptPassMock, tokenGenMock)
	assert.NoError(t, err)

	gotToken, gotSalt, err := authService.Login(ctx, login, password)
	assert.NoError(t, err)
	assert.Equal(t, token, gotToken)
	assert.Equal(t, salt, gotSalt)
}

func TestAuthService_Login_Error(t *testing.T) {
	ctx := context.Background()
	login := "testuser"
	password := "testpassword"

	userRepoMock := &mocks.UserRepository{}
	cryptPassMock := &mocks.ICryptPass{}
	tokenGenMock := &mocks.ITokenGenerator{}

	userRepoMock.EXPECT().Get(ctx, login).Return(entity.User{}, errors.New("user not found"))

	authService, err := NewAuthService(userRepoMock, cryptPassMock, tokenGenMock)
	assert.NoError(t, err)

	token, salt, err := authService.Login(ctx, login, password)
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Empty(t, salt)
	assert.Equal(t, "user testuser is not exist", err.Error())
}

func TestAuthService_Register_Success(t *testing.T) {
	ctx := context.Background()
	login := "newuser"
	password := "newpassword"
	hashedPassword := "hashedpassword"
	userID := uint64(2)
	salt := uuid.New().String()

	userRepoMock := &mocks.UserRepository{}
	cryptPassMock := &mocks.ICryptPass{}
	tokenGenMock := &mocks.ITokenGenerator{}

	cryptPassMock.EXPECT().HashPassword(password).Return(hashedPassword, nil)
	userRepoMock.EXPECT().Create(ctx, mock.Anything).Return(entity.User{
		ID:       userID,
		Login:    login,
		Password: hashedPassword,
		Salt:     salt,
	}, nil)

	authService, err := NewAuthService(userRepoMock, cryptPassMock, tokenGenMock)
	assert.NoError(t, err)

	user, err := authService.Register(ctx, login, password)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, login, user.Login)
	assert.Equal(t, hashedPassword, user.Password)
}

func TestAuthService_Register_Error(t *testing.T) {
	ctx := context.Background()
	login := "newuser"
	password := "newpassword"

	userRepoMock := &mocks.UserRepository{}
	cryptPassMock := &mocks.ICryptPass{}
	tokenGenMock := &mocks.ITokenGenerator{}

	cryptPassMock.EXPECT().HashPassword(password).Return("", errors.New("hash error"))

	authService, err := NewAuthService(userRepoMock, cryptPassMock, tokenGenMock)
	assert.NoError(t, err)

	user, err := authService.Register(ctx, login, password)
	assert.Error(t, err)
	assert.Equal(t, "hash error", err.Error())
	assert.Equal(t, entity.User{}, user)
}
