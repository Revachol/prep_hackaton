package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"auth-service/config"
	"auth-service/gen"
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/database"
	"auth-service/pkg/jwt"
	"auth-service/pkg/redis"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"

	_ "auth-service/docs" // Важно: импортируем сгенерированную документацию
)

// @title Auth Service API
// @version 1.0
// @description Microservice for user authentication and authorization with JWT tokens

// @contact.name API Support
// @contact.email support@auth-service.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем PostgreSQL
	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Инициализируем Redis
	redisClient, err := redis.NewRedisClient(cfg.RedisURL, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("✅ All databases connected successfully!")

	// Инициализируем JWT менеджер
	jwtManager := jwt.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)

	// Создаем репозитории
	tokenRepo := repository.NewTokenRepository(redisClient, cfg.JWTExpiration)
	userRepo := repository.NewUserRepository(db.DB)

	// Создаем сервисы
	authService := service.NewAuthService(jwtManager, tokenRepo, userRepo)

	// Создаем обработчики и middleware
	httpHandler := handlers.NewHTTPHandler(authService)
	grpcHandler := handlers.NewGRPCHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Запускаем gRPC сервер в горутине
	go startGRPCServer(cfg.GRPCPort, grpcHandler)

	// Настраиваем Gin
	router := gin.Default()

	// Добавляем CORS middleware ПЕРВЫМ
	router.Use(middleware.CORS())

	// Добавляем middleware для логирования
	router.Use(func(c *gin.Context) {
		log.Printf("HTTP Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Swagger документация
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	router.GET("/swagger/*any", swaggerHandler)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
			"services": gin.H{
				"postgres": "connected",
				"redis":    "connected",
			},
		})
	})

	// Группа публичных маршрутов
	public := router.Group("/api")
	{
		public.POST("/register", httpHandler.Register)
		public.POST("/login", httpHandler.Login)
	}

	// Группа защищенных маршрутов
	protected := router.Group("/api")
	protected.Use(authMiddleware.AuthRequired())
	{
		protected.POST("/logout", httpHandler.Logout)
		protected.GET("/me", httpHandler.Me)
	}

	// Ожидаем сигналов для graceful shutdown
	waitForShutdown()

	// Запускаем HTTP сервер
	log.Printf("🚀 Starting HTTP server on :%s", cfg.HTTPPort)
	log.Printf("🔐 gRPC server running on :%s", cfg.GRPCPort)
	log.Printf("📚 Swagger docs available at: http://localhost:%s/swagger/index.html", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// startGRPCServer запускает gRPC сервер
func startGRPCServer(port string, handler *handlers.GRPCHandler) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	gen.RegisterAuthServiceServer(grpcServer, handler)

	log.Printf("🚀 gRPC server listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

// waitForShutdown ожидает сигналов для graceful shutdown
func waitForShutdown() {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("🛑 Shutdown signal received...")
		os.Exit(0)
	}()
}
