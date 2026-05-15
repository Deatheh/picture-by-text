package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"picture-service/internal/ai"
	"picture-service/internal/config"
	"picture-service/internal/handler"
	"picture-service/internal/repository"
	"picture-service/internal/repository/db"
	"picture-service/internal/repository/minio"
	"picture-service/internal/service"

	pb "picturepb"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Конфиг
	envConf := config.NewEnvConfig()
	envConf.PrintConfig()

	// PostgreSQL
	dbRepo, err := db.NewDatabaseInstance(envConf)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer dbRepo.Close()

	// MinIO
	minioRepo, err := minio.NewMinioRepository(envConf)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	// GigaChat
	aiRepo, err := ai.InitAiRepository(envConf)
	if err != nil {
		log.Fatalf("Failed to init GigaChat: %v", err)
	}

	// Repository
	repo := &repository.Repository{
		DatabaseRepository: dbRepo,
		MinioRepository:    minioRepo,
	}

	// Service
	svc := service.NewService(repo, envConf, aiRepo)

	// gRPC Handler
	pictureHandler := handler.NewPictureHandler(envConf, svc)

	// gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", envConf.Application.Host, envConf.Application.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPictureServiceServer(grpcServer, pictureHandler)

	log.Printf("Picture Service starting on %s:%d", envConf.Application.Host, envConf.Application.Port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down picture service...")
	grpcServer.GracefulStop()
}
