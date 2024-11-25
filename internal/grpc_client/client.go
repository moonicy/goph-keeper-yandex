package grpc_client

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/moonicy/goph-keeper-yandex/crypt"
	pb "github.com/moonicy/goph-keeper-yandex/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
)

type Data struct {
	ID   uint64
	Data []byte
}

type Client struct {
	conn   *grpc.ClientConn
	client pb.GophKeeperClient
	token  string
}

func NewClient(target string) (*Client, error) {
	// Создаём CertPool и добавляем в него сертификат
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM([]byte(crypt.CaCert)) {
		log.Fatalf("Не удалось добавить сертификат в CertPool")
	}

	// Создаём креденшиалы TLS с использованием CertPool
	creds := credentials.NewClientTLSFromCert(certPool, "localhost")
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться: %v", err)
	}
	client := pb.NewGophKeeperClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Login(login string, password string) error {
	resp, err := c.client.LoginUser(context.Background(), &pb.LoginUserRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("ошибка авторизации: %v", err)
	}

	log.Println(resp.Message)

	c.token = resp.Token
	return nil
}

// Logout Очищаем токен
func (c *Client) Logout() {
	c.token = ""
}

func (c *Client) Register(login string, password string) (uint64, error) {
	resp, err := c.client.RegisterUser(context.Background(), &pb.RegisterUserRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return 0, fmt.Errorf("ошибка регистрации: %v", err)
	}

	log.Println(resp.Message)
	return resp.UserId, nil
}

func (c *Client) AddData(data []byte) error {
	resp, err := c.client.AddData(c.contextWithToken(), &pb.AddDataRequest{
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("ошибка добавления данных: %v", err)
	}
	log.Println(resp.Message)
	return nil
}

func (c *Client) GetData() ([]Data, error) {
	resp, err := c.client.GetData(c.contextWithToken(), &pb.GetDataRequest{})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных: %v", err)
	}
	return c.pbToData(resp.Data), nil
}

func (c *Client) UpdateData(id uint64, data []byte) error {
	resp, err := c.client.UpdateData(c.contextWithToken(), &pb.UpdateDataRequest{
		Id:   id,
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("ошибка обновления данных: %v", err)
	}
	log.Println(resp.Message)
	return nil
}

func (c *Client) RemoveData(id uint64) error {
	resp, err := c.client.RemoveData(c.contextWithToken(), &pb.RemoveDataRequest{
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("ошибка удаления данных: %v", err)
	}
	log.Println(resp.Message)
	return nil
}

// Функция для создания контекста с токеном
func (c *Client) contextWithToken() context.Context {
	ctx := context.Background()
	if c.token != "" {
		md := metadata.New(map[string]string{
			"authorization": "Bearer " + c.token,
		})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func (c *Client) pbToData(data []*pb.Data) []Data {
	dt := make([]Data, len(data))
	for i, d := range data {
		dt[i] = Data{
			ID:   d.Id,
			Data: d.Data,
		}
	}
	return dt
}
