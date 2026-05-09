package main

import (
	"api-gateway/internal/config"
	grpcclient "api-gateway/internal/grpc-client"
	"api-gateway/internal/handler"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	envConf := config.NewEnvConfig()
	envConf.PrintConfigWithHiddenSecrets()

	var userClient *grpcclient.UserClient
	userClient, err := grpcclient.NewUserClient(envConf.Services.UserServiceURL)
	if err != nil {
		log.Printf("WARNING: Failed to create user client: %v", err)
		log.Println("API Gateway will start but /auth/* endpoints will return errors")
		userClient = nil
	} else {
		defer userClient.Close()
	}

	handlers := handler.NewRouter(envConf, userClient)

	server := handlers.InitRoutes()

	go func() {
		log.Printf("API Gateway starting on %s", fmt.Sprintf("%v:%v", envConf.Application.Host, envConf.Application.Port))
		if err := server.Run(fmt.Sprintf("%v:%v", envConf.Application.Host, envConf.Application.Port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
}
