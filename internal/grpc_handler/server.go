package grpc_handler

import (
	"context"
	"errors"
	"github.com/moonicy/goph-keeper-yandex/internal/service"
	pb "github.com/moonicy/goph-keeper-yandex/proto"
)

// Server реализует интерфейс, сгенерированный из proto файла
type Server struct {
	pb.UnimplementedGophKeeperServer
	authService *service.AuthService
}

func NewServer(authService *service.AuthService) (*Server, error) {
	if authService == nil {
		return nil, errors.New("authService is nil")
	}
	return &Server{
		authService: authService,
	}, nil
}

// RegisterUser обрабатывает регистрацию пользователя
func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	userID, err := s.authService.Register(ctx, req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		UserId:  uint64(userID),
		Message: "User registered successfully",
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	token, err := s.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginUserResponse{
		Token:   token,
		Message: "Success login",
	}, nil
}
