package grpc_client

import (
	"context"
	pb "github.com/moonicy/goph-keeper-yandex/proto"
	"google.golang.org/grpc"
)

//go:generate mockery --output ./../../mocks --filename crypt_gen.go --outpkg mocks --name Crypt --with-expecter
type Crypt interface {
	Init(password string, salt string)
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(encryptedData []byte) ([]byte, error)
	Clean()
}

//go:generate mockery --output ./../../mocks --filename client_conn.go --outpkg mocks --name ClientConn --with-expecter
type ClientConn interface {
	Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error
	NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error)
	Close() error
}

//go:generate mockery --output ./../../mocks --filename goph_keeper_client.go --outpkg mocks --name GophKeeperClient --with-expecter
type GophKeeperClient interface {
	RegisterUser(ctx context.Context, in *pb.RegisterUserRequest, opts ...grpc.CallOption) (*pb.RegisterUserResponse, error)
	LoginUser(ctx context.Context, in *pb.LoginUserRequest, opts ...grpc.CallOption) (*pb.LoginUserResponse, error)
	AddData(ctx context.Context, in *pb.AddDataRequest, opts ...grpc.CallOption) (*pb.AddDataResponse, error)
	UpdateData(ctx context.Context, in *pb.UpdateDataRequest, opts ...grpc.CallOption) (*pb.UpdateDataResponse, error)
	GetData(ctx context.Context, in *pb.GetDataRequest, opts ...grpc.CallOption) (*pb.GetDataResponse, error)
	RemoveData(ctx context.Context, in *pb.RemoveDataRequest, opts ...grpc.CallOption) (*pb.RemoveDataResponse, error)
}
