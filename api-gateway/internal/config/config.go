package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Application struct {
	Host    string
	Port    int
	LogPath string
}

// Конфиги для подключения к сервисам
type Services struct {
	UserServiceURL    string
	PictureServiceURL string
	StorageServiceURL string
	ExportServiceURL  string
}

// PostgreSQL конфиг
type PostgreSQL struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

// MinIO конфиг
type MinIO struct {
	Endpoint         string
	AccessKey        string
	SecretKey        string
	UseSSL           bool
	BucketImages     string
	BucketExports    string
	ExternalEndpoint string
}

// Redis конфиг
type Redis struct {
	Host       string
	Port       int
	Password   string
	DB         int
	SessionDB  int
	CacheDB    int
	URL        string
	SessionTTL int // seconds
}

// JWT конфиг
type JWT struct {
	Secret     string
	AccessTTL  int // seconds
	RefreshTTL int // seconds
}

// Общий конфиг
type Config struct {
	Application Application
	Services    Services
	PostgreSQL  PostgreSQL
	MinIO       MinIO
	Redis       Redis
	JWT         JWT
}

func NewEnvConfig() *Config {
	return &Config{
		Application: Application{
			Host:    getEnv("API_GATEWAY_HOST", "0.0.0.0"),
			Port:    getEnvAsInt("API_GATEWAY_PORT", 8080),
			LogPath: getEnv("APP_LOG_PATH", "./logs"),
		},
		Services: Services{
			UserServiceURL:    getEnv("USER_SERVICE_URL", "user-service:50051"),
			PictureServiceURL: getEnv("PICTURE_SERVICE_URL", "picture-service:50052"),
			StorageServiceURL: getEnv("STORAGE_SERVICE_URL", "storage-service:50053"),
			ExportServiceURL:  getEnv("EXPORT_SERVICE_URL", "export-service:50054"),
		},
		PostgreSQL: PostgreSQL{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			DBName:   getEnv("POSTGRES_DB", "microservices"),
			SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
			URL:      getEnv("DATABASE_URL", ""),
		},
		MinIO: MinIO{
			Endpoint:         getEnv("MINIO_ENDPOINT", "minio:9000"),
			AccessKey:        getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:        getEnv("MINIO_SECRET_KEY", "minioadmin123"),
			UseSSL:           getEnvAsBool("MINIO_USE_SSL", false),
			BucketImages:     getEnv("MINIO_BUCKET_IMAGES", "images"),
			BucketExports:    getEnv("MINIO_BUCKET_EXPORTS", "exports"),
			ExternalEndpoint: getEnv("MINIO_EXTERNAL_ENDPOINT", "localhost:9000"),
		},
		Redis: Redis{
			Host:       getEnv("REDIS_HOST", "localhost"),
			Port:       getEnvAsInt("REDIS_PORT", 6379),
			Password:   getEnv("REDIS_PASSWORD", ""),
			DB:         getEnvAsInt("REDIS_DB", 0),
			SessionDB:  getEnvAsInt("REDIS_SESSION_DB", 1),
			CacheDB:    getEnvAsInt("REDIS_CACHE_DB", 2),
			URL:        getEnv("REDIS_URL", ""),
			SessionTTL: getEnvAsInt("SESSION_TTL", 86400),
		},
		JWT: JWT{
			Secret:     getEnv("JWT_SECRET", "secret-jwt-key"),
			AccessTTL:  getEnvAsInt("JWT_ACCESS_TTL", 900),
			RefreshTTL: getEnvAsInt("JWT_REFRESH_TTL", 604800),
		},
	}
}

// Хелперы для работы с env

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

// Маскирование секретов для вывода в лог
func (c *Config) maskSecret(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
}

func (c *Config) PrintConfigWithHiddenSecrets() {
	fmt.Println("========== CONFIGURATION ==========")

	fmt.Println("\n[Application]")
	fmt.Printf("  Host: %s\n", c.Application.Host)
	fmt.Printf("  Port: %d\n", c.Application.Port)
	fmt.Printf("  LogPath: %s\n", c.Application.LogPath)

	fmt.Println("\n[Services]")
	fmt.Printf("  User Service: %s\n", c.Services.UserServiceURL)
	fmt.Printf("  Picture Service: %s\n", c.Services.PictureServiceURL)
	fmt.Printf("  Storage Service: %s\n", c.Services.StorageServiceURL)
	fmt.Printf("  Export Service: %s\n", c.Services.ExportServiceURL)

	fmt.Println("\n[PostgreSQL]")
	fmt.Printf("  Host: %s\n", c.PostgreSQL.Host)
	fmt.Printf("  Port: %d\n", c.PostgreSQL.Port)
	fmt.Printf("  User: %s\n", c.PostgreSQL.User)
	fmt.Printf("  Password: %s\n", c.maskSecret(c.PostgreSQL.Password))
	fmt.Printf("  Database: %s\n", c.PostgreSQL.DBName)
	fmt.Printf("  SSL Mode: %s\n", c.PostgreSQL.SSLMode)

	fmt.Println("\n[MinIO]")
	fmt.Printf("  Endpoint: %s\n", c.MinIO.Endpoint)
	fmt.Printf("  Access Key: %s\n", c.MinIO.AccessKey)
	fmt.Printf("  Secret Key: %s\n", c.maskSecret(c.MinIO.SecretKey))
	fmt.Printf("  Use SSL: %t\n", c.MinIO.UseSSL)
	fmt.Printf("  Images Bucket: %s\n", c.MinIO.BucketImages)
	fmt.Printf("  Exports Bucket: %s\n", c.MinIO.BucketExports)

	fmt.Println("\n[Redis]")
	fmt.Printf("  Host: %s\n", c.Redis.Host)
	fmt.Printf("  Port: %d\n", c.Redis.Port)
	fmt.Printf("  DB: %d\n", c.Redis.DB)
	fmt.Printf("  Session DB: %d\n", c.Redis.SessionDB)
	fmt.Printf("  Cache DB: %d\n", c.Redis.CacheDB)
	fmt.Printf("  Session TTL: %d sec (%d hours)\n", c.Redis.SessionTTL, c.Redis.SessionTTL/3600)

	fmt.Println("\n[JWT]")
	fmt.Printf("  Secret: %s\n", c.maskSecret(c.JWT.Secret))
	fmt.Printf("  Access TTL: %d sec (%d min)\n", c.JWT.AccessTTL, c.JWT.AccessTTL/60)
	fmt.Printf("  Refresh TTL: %d sec (%d days)\n", c.JWT.RefreshTTL, c.JWT.RefreshTTL/86400)

	fmt.Println("===================================")
}
