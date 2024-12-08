package grpc_handler

import (
	"context"
	"errors"
	"github.com/moonicy/goph-keeper-yandex/internal/config"
	mocks2 "github.com/moonicy/goph-keeper-yandex/mocks"
	"testing"

	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	pb "github.com/moonicy/goph-keeper-yandex/proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_RegisterUser_Success(t *testing.T) {
	ctx := context.Background()

	authServiceMock := &mocks2.AuthService{}
	dataRepoMock := &mocks2.DataRepository{}

	expectedUser := entity.User{
		ID:    uint64(1),
		Login: "testuser",
	}
	authServiceMock.EXPECT().Register(ctx, "testuser", "testpassword").
		Return(expectedUser, nil)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.RegisterUserRequest{
		Login:    "testuser",
		Password: "testpassword",
	}
	resp, err := server.RegisterUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(1), resp.UserId)
	assert.Equal(t, "User registered successfully", resp.Message)

	authServiceMock.AssertExpectations(t)
}

func TestServer_RegisterUser_Error(t *testing.T) {
	ctx := context.Background()

	authServiceMock := &mocks2.AuthService{}
	dataRepoMock := &mocks2.DataRepository{}

	expectedError := errors.New("registration failed")
	authServiceMock.EXPECT().Register(ctx, "testuser", "testpassword").
		Return(entity.User{}, expectedError)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.RegisterUserRequest{
		Login:    "testuser",
		Password: "testpassword",
	}
	resp, err := server.RegisterUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)

	authServiceMock.AssertExpectations(t)
}

func TestServer_LoginUser_Success(t *testing.T) {
	ctx := context.Background()

	authServiceMock := &mocks2.AuthService{}
	dataRepoMock := &mocks2.DataRepository{}

	expectedToken := "token123"
	expectedSalt := "salt123"
	authServiceMock.EXPECT().Login(ctx, "testuser", "testpassword").
		Return(expectedToken, expectedSalt, nil)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.LoginUserRequest{
		Login:    "testuser",
		Password: "testpassword",
	}
	resp, err := server.LoginUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedToken, resp.Token)
	assert.Equal(t, "Success login", resp.Message)
	assert.Equal(t, expectedSalt, resp.Salt)

	authServiceMock.AssertExpectations(t)
}

func TestServer_LoginUser_Error(t *testing.T) {
	ctx := context.Background()

	authServiceMock := &mocks2.AuthService{}
	dataRepoMock := &mocks2.DataRepository{}

	expectedError := errors.New("login failed")
	authServiceMock.EXPECT().Login(ctx, "testuser", "testpassword").
		Return("", "", expectedError)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.LoginUserRequest{
		Login:    "testuser",
		Password: "testpassword",
	}
	resp, err := server.LoginUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)

	authServiceMock.AssertExpectations(t)
}

func TestServer_AddData_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), "user_id", uint64(1))
	testdata := []byte("testdata")

	dataRepoMock := &mocks2.DataRepository{}
	authServiceMock := &mocks2.AuthService{}

	dataRepoMock.EXPECT().AddData(ctx, entity.Data{
		UserID: 1,
		Data:   testdata,
	}).Return(nil)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.AddDataRequest{Data: testdata}
	resp, err := server.AddData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Data added successfully", resp.Message)

	dataRepoMock.AssertExpectations(t)
}

func TestServer_AddData_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), "user_id", uint64(1))
	testdata := []byte("testdata")

	dataRepoMock := &mocks2.DataRepository{}
	authServiceMock := &mocks2.AuthService{}

	expectedError := errors.New("failed to add data")
	dataRepoMock.EXPECT().AddData(ctx, entity.Data{
		UserID: 1,
		Data:   testdata,
	}).Return(expectedError)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.AddDataRequest{Data: testdata}
	resp, err := server.AddData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)

	dataRepoMock.AssertExpectations(t)
}

func TestServer_UpdateData_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), "user_id", uint64(1))
	updatedData := []byte("updatedData")

	dataRepoMock := &mocks2.DataRepository{}
	authServiceMock := &mocks2.AuthService{}

	dataRepoMock.EXPECT().UpdateData(ctx, entity.Data{
		UserID: 1,
		ID:     1,
		Data:   updatedData,
	}).Return(nil)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.UpdateDataRequest{Id: 1, Data: updatedData}
	resp, err := server.UpdateData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Data updated successfully", resp.Message)

	dataRepoMock.AssertExpectations(t)
}

func TestServer_UpdateData_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), "user_id", uint64(1))
	updatedData := []byte("updatedData")

	dataRepoMock := &mocks2.DataRepository{}
	authServiceMock := &mocks2.AuthService{}

	expectedError := errors.New("update failed")
	dataRepoMock.EXPECT().UpdateData(ctx, entity.Data{
		UserID: 1,
		ID:     1,
		Data:   updatedData,
	}).Return(expectedError)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.UpdateDataRequest{Id: 1, Data: updatedData}
	resp, err := server.UpdateData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)

	dataRepoMock.AssertExpectations(t)
}

func TestServer_RemoveData_Success(t *testing.T) {
	ctx := context.Background()

	dataRepoMock := &mocks2.DataRepository{}
	authServiceMock := &mocks2.AuthService{}

	dataRepoMock.EXPECT().RemoveData(ctx, uint64(1)).Return(nil)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.RemoveDataRequest{Id: 1}
	resp, err := server.RemoveData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Data removed successfully", resp.Message)

	dataRepoMock.AssertExpectations(t)
}

func TestServer_RemoveData_Error(t *testing.T) {
	ctx := context.Background()

	dataRepoMock := &mocks2.DataRepository{}
	authServiceMock := &mocks2.AuthService{}

	expectedError := errors.New("failed to remove data")
	dataRepoMock.EXPECT().RemoveData(ctx, uint64(1)).Return(expectedError)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.RemoveDataRequest{Id: 1}
	resp, err := server.RemoveData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)

	dataRepoMock.AssertExpectations(t)
}

func TestServer_GetData_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), "user_id", uint64(1))
	data1 := []byte("data1")
	data2 := []byte("data2")

	dataRepoMock := &mocks2.DataRepository{}
	authServiceMock := &mocks2.AuthService{}

	expectedData := []entity.Data{
		{ID: 1, UserID: 1, Data: data1},
		{ID: 2, UserID: 1, Data: data2},
	}

	dataRepoMock.EXPECT().GetData(ctx, uint64(1)).Return(expectedData, nil)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.GetDataRequest{}
	resp, err := server.GetData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Data, len(expectedData))

	for i, d := range resp.Data {
		assert.Equal(t, expectedData[i].ID, d.Id)
		assert.Equal(t, expectedData[i].Data, d.Data)
	}

	dataRepoMock.AssertExpectations(t)
}

func TestServer_GetData_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), "user_id", uint64(1))

	dataRepoMock := &mocks2.DataRepository{}
	authServiceMock := &mocks2.AuthService{}

	expectedError := errors.New("failed to get data")
	dataRepoMock.EXPECT().GetData(ctx, uint64(1)).Return(nil, expectedError)

	server, err := NewServer(authServiceMock, dataRepoMock, config.ServerConfig{
		CryptoKey: "testkey",
		CryptoCrt: "testcrt",
	})
	assert.NoError(t, err)

	req := &pb.GetDataRequest{}
	resp, err := server.GetData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)

	dataRepoMock.AssertExpectations(t)
}
