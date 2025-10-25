package main

import (
	"log"
	"os"
	"time"

	"backend/internal/auth"
	"backend/internal/modules/user/delivery"
	"backend/internal/modules/user/repository"
	"backend/internal/modules/user/usecase"
	"backend/pkg/db"
	"github.com/gin-gonic/gin"
)

func main() {
	// -------------------------
	// Настройка БД
	// -------------------------
	dsn := "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
	pg := db.NewPostgresDB(dsn)

	// -------------------------
	// AuthManager (JWT + Redis)
	// -------------------------
	secret := "supersecret" // в проде лучше из env
	tokenTTL := time.Hour * 24
	// var redisClient interface{} = nil // пока nil для хакатона
	authManager := auth.NewSessionJWTManager(secret, tokenTTL, nil)

	// -------------------------
	// Репозиторий и сервис
	// -------------------------
	userRepo := repository.NewPostgresUserRepository(pg)
	userService := usecase.NewUserService(userRepo, authManager)

	// -------------------------
	// Gin и HTTP-хендлеры
	// -------------------------
	r := gin.Default()
	delivery.NewUserHandler(r, userService, authManager)

	// -------------------------
	// Запуск сервера
	// -------------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
