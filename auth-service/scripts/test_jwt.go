package main

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	"auth-service/pkg/jwt"
// )

// func main() {
// 	// Тестируем JWT менеджер
// 	jwtManager := jwt.NewJWTManager("test-secret-key", 24*time.Hour)

// 	// Генерируем токен
// 	userID := 1
// 	version := 1
// 	token, err := jwtManager.GenerateToken(userID, version)
// 	if err != nil {
// 		log.Fatal("Error generating token:", err)
// 	}

// 	fmt.Printf("Generated token: %s\n\n", token)

// 	// Верифицируем токен
// 	claims, err := jwtManager.VerifyToken(token)
// 	if err != nil {
// 		log.Fatal("Error verifying token:", err)
// 	}

// 	fmt.Printf("Token verified successfully!\n")
// 	fmt.Printf("UserID: %d\n", claims.UserID)
// 	fmt.Printf("Version: %d\n", claims.Version)

// 	// Пробуем извлечь claims (для отладки)
// 	extractedClaims, err := jwtManager.ExtractClaims(token)
// 	if err != nil {
// 		log.Fatal("Error extracting claims:", err)
// 	}

// 	fmt.Printf("\nExtracted claims:\n")
// 	fmt.Printf("UserID: %d\n", extractedClaims.UserID)
// 	fmt.Printf("Version: %d\n", extractedClaims.Version)
// 	fmt.Printf("Expires: %s\n", extractedClaims.ExpiresAt.Time)
// }
