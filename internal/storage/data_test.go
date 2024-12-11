package storage

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewDataRepository_Success(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo, err := NewDataRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func TestNewDataRepository_NilDB(t *testing.T) {
	repo, err := NewDataRepository(nil)
	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Equal(t, "db is nil", err.Error())
}

func TestDataRepository_AddData_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	assert.NoError(t, err)

	repo, err := NewDataRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	data := entity.Data{
		UserID: 1,
		Data:   []byte("test data"),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "data" ("user_id","data") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(data.UserID, data.Data).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.AddData(ctx, data)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDataRepository_AddData_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	assert.NoError(t, err)

	repo, err := NewDataRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	data := entity.Data{
		UserID: 1,
		Data:   []byte("test data"),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "data" ("user_id","data") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(data.UserID, data.Data).
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	err = repo.AddData(ctx, data)
	assert.Error(t, err)
	assert.Equal(t, "insert error", err.Error())

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDataRepository_RemoveData_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	assert.NoError(t, err)

	repo, err := NewDataRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	id := uint64(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "data" WHERE "data"."id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.RemoveData(ctx, id)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDataRepository_RemoveData_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	assert.NoError(t, err)

	repo, err := NewDataRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	id := uint64(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "data" WHERE "data"."id" = $1`)).
		WithArgs(id).
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	err = repo.RemoveData(ctx, id)
	assert.Error(t, err)
	assert.Equal(t, "delete error", err.Error())

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDataRepository_GetData_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	assert.NoError(t, err)

	repo, err := NewDataRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	userID := uint64(1)

	dataRows := sqlmock.NewRows([]string{"id", "user_id", "data"}).
		AddRow(1, userID, []byte("test data 1")).
		AddRow(2, userID, []byte("test data 2"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "data" WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(dataRows)

	result, err := repo.GetData(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, uint64(1), result[0].ID)
	assert.Equal(t, userID, result[0].UserID)
	assert.Equal(t, []byte("test data 1"), result[0].Data)
	assert.Equal(t, uint64(2), result[1].ID)
	assert.Equal(t, userID, result[1].UserID)
	assert.Equal(t, []byte("test data 2"), result[1].Data)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDataRepository_GetData_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	assert.NoError(t, err)

	repo, err := NewDataRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	userID := uint64(1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "data" WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnError(errors.New("select error"))

	result, err := repo.GetData(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, "select error", err.Error())
	assert.Empty(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDataRepository_UpdateData_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	assert.NoError(t, err)

	repo, err := NewDataRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	data := entity.Data{
		ID:     1,
		UserID: 1,
		Data:   []byte("updated data"),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "data" SET "user_id"=$1,"data"=$2 WHERE "id" = $3`)).
		WithArgs(data.UserID, data.Data, data.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.UpdateData(ctx, data)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDataRepository_UpdateData_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	assert.NoError(t, err)

	repo, err := NewDataRepository(gormDB)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	ctx := context.Background()

	data := entity.Data{
		ID:     1,
		UserID: 1,
		Data:   []byte("updated data"),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "data" SET "user_id"=$1,"data"=$2 WHERE "id" = $3`)).
		WithArgs(data.UserID, data.Data, data.ID).
		WillReturnError(errors.New("update error"))
	mock.ExpectRollback()

	err = repo.UpdateData(ctx, data)
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
