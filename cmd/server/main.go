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
	jwtKey := "popa"

	db, err := storage.NewDB(config.ServerConfig{Database: "host=localhost port=5432 user=mila dbname=goph_keeper password=qwerty sslmode=disable"})
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
	authInter := interceptor.NewAuthInterceptor(jwtKey)
	cryptPass := service.NewCryptPass()
	gen, err := service.NewTokenGenerator(jwtKey)
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

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("Server is running on port 8080...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
