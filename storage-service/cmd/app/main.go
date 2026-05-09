package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"storage-service/internal/config"
	"storage-service/internal/handler"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb "storagepb"
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

	storageHandler := handler.NewStorageHandler(envConf)

	pb.RegisterStorageServiceServer(grpcServer, storageHandler)

	log.Printf("storage-service starting on port %d", envConf.Application.Port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down storage-service...")
	grpcServer.GracefulStop()
}
