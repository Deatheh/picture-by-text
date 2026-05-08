package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/handler"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	envConf := config.NewEnvConfig()
	envConf.PrintConfigWithHiddenSecrets()

	handlers := handler.NewRouter(envConf)

	if err := handlers.InitRoutes().Run(fmt.Sprintf(":%v", envConf.Application.Port)); err != nil {
		log.Fatal(fmt.Errorf("server run error: %w", err))
	}
}
