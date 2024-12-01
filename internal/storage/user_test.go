package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewUserRepository_Success(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo, err := NewUserRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func TestNewUserRepository_NilDB(t *testing.T) {
	repo, err := NewUserRepository(nil)
	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Equal(t, "NewUserRepository: db is nil", err.Error())
}

func TestUserRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo, err := NewUserRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	user := entity.User{
		Login:    "testuser",
		Password: "hashedpassword",
		Salt:     "randomsalt",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" .* RETURNING "id"`).
		WithArgs(user.Login, user.Password, user.Salt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	result, err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.ID)
	assert.Equal(t, user.Login, result.Login)
	assert.Equal(t, user.Password, result.Password)
	assert.Equal(t, user.Salt, result.Salt)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo, err := NewUserRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	user := entity.User{
		Login:    "testuser",
		Password: "hashedpassword",
		Salt:     "randomsalt",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" .* RETURNING "id"`).
		WithArgs(user.Login, user.Password, user.Salt).
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	result, err := repo.Create(ctx, user)
	assert.Error(t, err)
	assert.Equal(t, "insert error", err.Error())
	assert.Equal(t, entity.User{}, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_Get_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo, err := NewUserRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	login := "testuser"
	user := entity.User{
		ID:       1,
		Login:    login,
		Password: "hashedpassword",
		Salt:     "randomsalt",
	}

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE login = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(login, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "salt"}).
			AddRow(user.ID, user.Login, user.Password, user.Salt))

	result, err := repo.Get(ctx, login)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo, err := NewUserRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	login := "nonexistentuser"

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE login = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(login, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "salt"}))

	result, err := repo.Get(ctx, login)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.Equal(t, entity.User{}, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_Get_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo, err := NewUserRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	login := "testuser"

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE login = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(login, 1).
		WillReturnError(errors.New("query error"))

	result, err := repo.Get(ctx, login)
	assert.Error(t, err)
	assert.Equal(t, "query error", err.Error())
	assert.Equal(t, entity.User{}, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
