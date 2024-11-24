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
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.dataRepo.AddData(ctx, entity.Data{
		UserID: userID,
		Data:   req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddDataResponse{Message: "Data added successfully"}, nil
}

func (s *Server) UpdateData(ctx context.Context, req *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.dataRepo.UpdateData(ctx, entity.Data{
		UserID: userID,
		ID:     req.Id,
		Data:   req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateDataResponse{Message: "Data updated successfully"}, nil
}

func (s *Server) RemoveData(ctx context.Context, req *pb.RemoveDataRequest) (*pb.RemoveDataResponse, error) {
	err := s.dataRepo.RemoveData(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveDataResponse{Message: "Data removed successfully"}, nil
}

func (s *Server) GetData(ctx context.Context, _ *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	data, err := s.dataRepo.GetData(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &pb.GetDataResponse{
		Data: s.entityDataToPB(data),
	}, nil
}

func (s *Server) getUserID(ctx context.Context) (uint64, error) {
	ctxValue := ctx.Value("user_id")
	var userID uint64

	switch v := ctxValue.(type) {
	case uint64:
		userID = v
	default:
		return 0, errors.New("user id is invalid")
	}
	return userID, nil
}

func (s *Server) entityDataToPB(data []entity.Data) []*pb.Data {
	var pbData []*pb.Data
	for _, d := range data {
		pbData = append(pbData, &pb.Data{
			Id:   d.ID,
			Data: d.Data,
		})
	}
	return pbData
}