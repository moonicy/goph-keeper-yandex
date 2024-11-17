package grpc_handler

import (
	"context"
	"errors"
	"github.com/moonicy/goph-keeper-yandex/internal/entity"
	"github.com/moonicy/goph-keeper-yandex/internal/service"
	"github.com/moonicy/goph-keeper-yandex/internal/storage"
	pb "github.com/moonicy/goph-keeper-yandex/proto"
)

// Server реализует интерфейс, сгенерированный из proto файла
type Server struct {
	pb.UnimplementedGophKeeperServer
	authService *service.AuthService
	dataRepo    *storage.DataRepository
}

func NewServer(authService *service.AuthService, dataRepo *storage.DataRepository) (*Server, error) {
	if authService == nil {
		return nil, errors.New("authService is nil")
	}
	if dataRepo == nil {
		return nil, errors.New("dataRepo is nil")
	}
	return &Server{
		authService: authService,
		dataRepo:    dataRepo,
	}, nil
}

// RegisterUser обрабатывает регистрацию пользователя
func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	userID, err := s.authService.Register(ctx, req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		UserId:  userID,
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

func (s *Server) AddData(ctx context.Context, req *pb.AddDataRequest) (*pb.AddDataResponse, error) {
	ctxValue := ctx.Value("user_id")
	var userID uint64

	switch v := ctxValue.(type) {
	case uint64:
		userID = v
	default:
		return nil, errors.New("user id is invalid")
	}

	err := s.dataRepo.AddData(ctx, entity.Data{
		UserID: userID,
		Data:   req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddDataResponse{Message: "Data added successfully"}, nil
}
