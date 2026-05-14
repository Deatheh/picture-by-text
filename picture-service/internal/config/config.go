package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Application struct {
	Host string
	Port int
}

type PostgreSQL struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type MinIO struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	BucketName string
}

type GigaChat struct {
	AuthKey string
}

type Config struct {
	Application Application
	PostgreSQL  PostgreSQL
	MinIO       MinIO
	GigaChat    GigaChat
}

func NewEnvConfig() *Config {
	return &Config{
		Application: Application{
			Host: getEnv("PICTURE_SERVICE_HOST", "picture-service"),
			Port: getEnvAsInt("PICTURE_SERVICE_PORT", 50052),
		},
		PostgreSQL: PostgreSQL{
			Host:     getEnv("POSTGRES_HOST", "postgres"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			DBName:   getEnv("POSTGRES_DB", "microservices"),
			SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
		},
		MinIO: MinIO{
			Endpoint:   getEnv("MINIO_ENDPOINT", "minio:9000"),
			AccessKey:  getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:  getEnv("MINIO_SECRET_KEY", "minioadmin123"),
			UseSSL:     getEnvAsBool("MINIO_USE_SSL", false),
			BucketName: getEnv("MINIO_BUCKET_IMAGES", "images"),
		},
		GigaChat: GigaChat{
			AuthKey: getEnv("GIGACHAT_AUTH_KEY", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}

func maskSecret(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
}

func (c *Config) PrintConfig() {
	fmt.Println("========== PICTURE SERVICE CONFIGURATION ==========")

	fmt.Println("\n[Application]")
	fmt.Printf("  Host: %s\n", c.Application.Host)
	fmt.Printf("  Port: %d\n", c.Application.Port)

	fmt.Println("\n[PostgreSQL]")
	fmt.Printf("  Host: %s\n", c.PostgreSQL.Host)
	fmt.Printf("  Port: %d\n", c.PostgreSQL.Port)
	fmt.Printf("  User: %s\n", c.PostgreSQL.User)
	fmt.Printf("  Password: %s\n", maskSecret(c.PostgreSQL.Password))
	fmt.Printf("  Database: %s\n", c.PostgreSQL.DBName)
	fmt.Printf("  SSL Mode: %s\n", c.PostgreSQL.SSLMode)

	fmt.Println("\n[MinIO]")
	fmt.Printf("  Endpoint: %s\n", c.MinIO.Endpoint)
	fmt.Printf("  Access Key: %s\n", c.MinIO.AccessKey)
	fmt.Printf("  Secret Key: %s\n", maskSecret(c.MinIO.SecretKey))
	fmt.Printf("  Use SSL: %t\n", c.MinIO.UseSSL)
	fmt.Printf("  Bucket: %s\n", c.MinIO.BucketName)

	fmt.Println("\n[GigaChat]")
	fmt.Printf("  Auth Key: %s\n", maskSecret(c.GigaChat.AuthKey))

	fmt.Println("====================================================")
}
