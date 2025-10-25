package config

import (
	"os"
	"regexp"
	"strconv"
	"time"
)

type Config struct {
	// HTTP Server
	HTTPPort string

	// gRPC Server
	GRPCPort string

	// Database
	DatabaseURL string

	// Redis
	RedisURL      string
	RedisPassword string
	RedisDB       int

	// JWT
	JWTSecret     string
	JWTExpiration time.Duration

	// Bcrypt
	BcryptCost int
}

func Load() *Config {
	cfg := &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		GRPCPort:      getEnv("GRPC_PORT", "50051"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://auth_user:auth_password@localhost:5434/auth_db?sslmode=disable"),
		RedisURL:      getEnv("REDIS_URL", "localhost:6378"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		BcryptCost:    getEnvAsInt("BCRYPT_COST", 10),
	}

	// Логируем настройки Redis (без пароля)
	println("Configuration loaded:")
	println("  Database URL:", maskPassword(cfg.DatabaseURL))
	println("  Redis URL:", cfg.RedisURL)
	println("  Bcrypt Cost:", cfg.BcryptCost)

	return cfg
}

func maskPassword(url string) string {
	// Простая маскировка пароля в URL
	// postgres://user:password@host/db -> postgres://user:****@host/db
	return regexp.MustCompile(`://([^:]+):[^@]+@`).ReplaceAllString(url, "://$1:****@")
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

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
