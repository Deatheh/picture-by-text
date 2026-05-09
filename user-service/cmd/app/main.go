package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"user-service/internal/config"
	"user-service/internal/handler"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb "userpb"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	envConf := config.NewEnvConfig()
	envConf.PrintConfigWithHiddenSecrets()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", envConf.Application.Host, envConf.Application.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	userHandler := handler.NewUserHandler()

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	log.Printf("user-service starting on port %d", envConf.Application.Port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down user-service...")
	grpcServer.GracefulStop()
}
