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

// Конфиги для подключения к сервисам
type Services struct {
	StorageServiceURL string
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
	Redis       Redis
	JWT         JWT
}

func NewEnvConfig() *Config {
	return &Config{
		Application: Application{
			Host:    getEnv("USER_SERVICE_HOST", "user-service"),
			Port:    getEnvAsInt("USER_SERVICE_PORT", 50051),
			LogPath: getEnv("APP_LOG_PATH", "./logs"),
		},
		Services: Services{
			StorageServiceURL: getEnv("STORAGE_SERVICE_URL", "storage-service:50053"),
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
	fmt.Printf("  Storage Service: %s\n", c.Services.StorageServiceURL)

	fmt.Println("\n[PostgreSQL]")
	fmt.Printf("  Host: %s\n", c.PostgreSQL.Host)
	fmt.Printf("  Port: %d\n", c.PostgreSQL.Port)
	fmt.Printf("  User: %s\n", c.PostgreSQL.User)
	fmt.Printf("  Password: %s\n", c.maskSecret(c.PostgreSQL.Password))
	fmt.Printf("  Database: %s\n", c.PostgreSQL.DBName)
	fmt.Printf("  SSL Mode: %s\n", c.PostgreSQL.SSLMode)

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
