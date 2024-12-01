package grpc_client

import (
	"context"
	"errors"
	"testing"

	"github.com/moonicy/goph-keeper-yandex/mocks"
	pb "github.com/moonicy/goph-keeper-yandex/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

func TestClient_Login_Success(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
	}

	login := "testuser"
	password := "testpassword"
	token := "testtoken"
	salt := "testsalt"

	clientMock.EXPECT().
		LoginUser(mock.Anything, &pb.LoginUserRequest{
			Login:    login,
			Password: password,
		}).
		Return(&pb.LoginUserResponse{
			Token:   token,
			Message: "Success login",
			Salt:    salt,
		}, nil)

	cryptMock.EXPECT().Init(password, salt)

	err := client.Login(login, password)

	assert.NoError(t, err)
	assert.Equal(t, token, client.token)

	clientMock.AssertExpectations(t)
	cryptMock.AssertExpectations(t)
}

func TestClient_Login_Error(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
	}

	login := "testuser"
	password := "testpassword"

	expectedError := errors.New("authorization error")

	clientMock.EXPECT().
		LoginUser(mock.Anything, &pb.LoginUserRequest{
			Login:    login,
			Password: password,
		}).
		Return(nil, expectedError)

	err := client.Login(login, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка авторизации")
	assert.Empty(t, client.token)

	clientMock.AssertExpectations(t)
}

func TestClient_Logout(t *testing.T) {
	cryptMock := &mocks.Crypt{}
	client := &Client{
		token: "testtoken",
		cpt:   cryptMock,
	}

	cryptMock.EXPECT().Clean()

	client.Logout()

	assert.Empty(t, client.token)
	cryptMock.AssertExpectations(t)
}

func TestClient_Register_Success(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
	}

	login := "newuser"
	password := "newpassword"
	userID := uint64(1)

	clientMock.EXPECT().
		RegisterUser(mock.Anything, &pb.RegisterUserRequest{
			Login:    login,
			Password: password,
		}).
		Return(&pb.RegisterUserResponse{
			UserId:  userID,
			Message: "User registered successfully",
		}, nil)

	id, err := client.Register(login, password)

	assert.NoError(t, err)
	assert.Equal(t, userID, id)

	clientMock.AssertExpectations(t)
}

func TestClient_Register_Error(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
	}

	login := "newuser"
	password := "newpassword"

	expectedError := errors.New("registration error")

	clientMock.EXPECT().
		RegisterUser(mock.Anything, &pb.RegisterUserRequest{
			Login:    login,
			Password: password,
		}).
		Return(nil, expectedError)

	id, err := client.Register(login, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка регистрации")
	assert.Zero(t, id)

	clientMock.AssertExpectations(t)
}

func TestClient_AddData_Success(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	data := []byte("testdata")
	encryptedData := []byte("encrypteddata")

	cryptMock.EXPECT().Encrypt(data).Return(encryptedData, nil)

	clientMock.EXPECT().
		AddData(mock.Anything, &pb.AddDataRequest{
			Data: encryptedData,
		}).
		Return(&pb.AddDataResponse{
			Message: "Data added successfully",
		}, nil)

	err := client.AddData(data)

	assert.NoError(t, err)

	cryptMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
}

func TestClient_AddData_Error_Encrypt(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	data := []byte("testdata")
	expectedError := errors.New("encryption error")

	cryptMock.EXPECT().Encrypt(data).Return(nil, expectedError)

	err := client.AddData(data)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка шифрования данных")

	cryptMock.AssertExpectations(t)
	clientMock.AssertNotCalled(t, "AddData", mock.Anything, mock.Anything)
}

func TestClient_AddData_Error_AddData(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	data := []byte("testdata")
	encryptedData := []byte("encrypteddata")
	expectedError := errors.New("add data error")

	cryptMock.EXPECT().Encrypt(data).Return(encryptedData, nil)

	clientMock.EXPECT().
		AddData(mock.Anything, &pb.AddDataRequest{
			Data: encryptedData,
		}).
		Return(nil, expectedError)

	err := client.AddData(data)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка добавления данных")

	cryptMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
}

func TestClient_GetData_Success(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	encryptedData := []byte("encrypteddata")
	decryptedData := []byte("testdata")

	pbData := []*pb.Data{
		{Id: 1, Data: encryptedData},
	}

	clientMock.EXPECT().
		GetData(mock.Anything, &pb.GetDataRequest{}).
		Return(&pb.GetDataResponse{
			Data: pbData,
		}, nil)

	cryptMock.EXPECT().
		Decrypt(encryptedData).
		Return(decryptedData, nil)

	data, err := client.GetData()

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, uint64(1), data[0].ID)
	assert.Equal(t, decryptedData, data[0].Data)

	clientMock.AssertExpectations(t)
	cryptMock.AssertExpectations(t)
}

func TestClient_GetData_Error_GetData(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	expectedError := errors.New("get data error")

	clientMock.EXPECT().
		GetData(mock.Anything, &pb.GetDataRequest{}).
		Return(nil, expectedError)

	data, err := client.GetData()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка получения данных")
	assert.Nil(t, data)

	clientMock.AssertExpectations(t)
	cryptMock.AssertNotCalled(t, "Decrypt", mock.Anything)
}

func TestClient_GetData_Error_Decrypt(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	encryptedData := []byte("encrypteddata")
	pbData := []*pb.Data{
		{Id: 1, Data: encryptedData},
	}

	expectedError := errors.New("decrypt error")

	clientMock.EXPECT().
		GetData(mock.Anything, &pb.GetDataRequest{}).
		Return(&pb.GetDataResponse{
			Data: pbData,
		}, nil)

	cryptMock.EXPECT().
		Decrypt(encryptedData).
		Return(nil, expectedError)

	data, err := client.GetData()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка расшифровки данных")
	assert.Nil(t, data)

	clientMock.AssertExpectations(t)
	cryptMock.AssertExpectations(t)
}

func TestClient_UpdateData_Success(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	id := uint64(1)
	data := []byte("updated data")
	encryptedData := []byte("encrypted updated data")

	cryptMock.EXPECT().Encrypt(data).Return(encryptedData, nil)

	clientMock.EXPECT().
		UpdateData(mock.Anything, &pb.UpdateDataRequest{
			Id:   id,
			Data: encryptedData,
		}).
		Return(&pb.UpdateDataResponse{
			Message: "Data updated successfully",
		}, nil)

	err := client.UpdateData(id, data)

	assert.NoError(t, err)

	cryptMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
}

func TestClient_UpdateData_Error(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	id := uint64(1)
	data := []byte("updated data")
	encryptedData := []byte("encrypted updated data")
	expectedError := errors.New("update data error")

	cryptMock.EXPECT().Encrypt(data).Return(encryptedData, nil)

	clientMock.EXPECT().
		UpdateData(mock.Anything, &pb.UpdateDataRequest{
			Id:   id,
			Data: encryptedData,
		}).
		Return(nil, expectedError)

	err := client.UpdateData(id, data)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка обновления данных")

	cryptMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
}

func TestClient_RemoveData_Success(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	id := uint64(1)

	clientMock.EXPECT().
		RemoveData(mock.Anything, &pb.RemoveDataRequest{
			Id: id,
		}).
		Return(&pb.RemoveDataResponse{
			Message: "Data removed successfully",
		}, nil)

	err := client.RemoveData(id)

	assert.NoError(t, err)

	clientMock.AssertExpectations(t)
}

func TestClient_RemoveData_Error(t *testing.T) {
	connMock := &mocks.ClientConn{}
	clientMock := &mocks.GophKeeperClient{}
	cryptMock := &mocks.Crypt{}

	client := &Client{
		conn:   connMock,
		client: clientMock,
		cpt:    cryptMock,
		token:  "testtoken",
	}

	id := uint64(1)
	expectedError := errors.New("remove data error")

	clientMock.EXPECT().
		RemoveData(mock.Anything, &pb.RemoveDataRequest{
			Id: id,
		}).
		Return(nil, expectedError)

	err := client.RemoveData(id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка удаления данных")

	clientMock.AssertExpectations(t)
}

func TestClient_Close(t *testing.T) {
	connMock := &mocks.ClientConn{}
	client := &Client{
		conn: connMock,
	}

	connMock.EXPECT().Close().Return(nil)

	client.Close()

	connMock.AssertExpectations(t)
}

func TestClient_contextWithToken(t *testing.T) {
	client := &Client{
		token: "testtoken",
	}

	ctx := client.contextWithToken()
	md, ok := metadata.FromOutgoingContext(ctx)
	assert.True(t, ok)
	assert.Contains(t, md.Get("authorization"), "Bearer testtoken")
}

func TestClient_contextWithToken_NoToken(t *testing.T) {
	client := &Client{}

	ctx := client.contextWithToken()
	_, ok := metadata.FromOutgoingContext(ctx)
	assert.False(t, ok)
	assert.Equal(t, context.Background(), ctx)
}
