package main

import (
	"fmt"
	"github.com/moonicy/goph-keeper-yandex/internal/config"
	"github.com/moonicy/goph-keeper-yandex/internal/grpc_handler"
	"github.com/moonicy/goph-keeper-yandex/internal/interceptor"
	"github.com/moonicy/goph-keeper-yandex/internal/service"
	"github.com/moonicy/goph-keeper-yandex/internal/storage"
	pb "github.com/moonicy/goph-keeper-yandex/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.NewServerConfig()

	db, err := storage.NewDB(cfg.DatabaseDsn)
	if err != nil {
		log.Fatal(err)
	}
	userRepo, err := storage.NewUserRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	dataRepo, err := storage.NewDataRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	authInter := interceptor.NewAuthInterceptor(cfg.JwtKey)
	cryptPass := service.NewCryptPass()
	gen, err := service.NewTokenGenerator(cfg.JwtKey)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := service.NewAuthService(userRepo, cryptPass, gen)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer, err := grpc_handler.NewServer(auth, dataRepo)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(authInter.Unary()))
	pb.RegisterGophKeeperServer(server, grpcServer)

	listener, err := net.Listen("tcp", cfg.Host+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Printf("Server is running on port %s", cfg.Port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
