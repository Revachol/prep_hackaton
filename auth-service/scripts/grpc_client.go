package main

import (
	"context"
	"log"
	"time"

	"auth-service/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Подключаемся к gRPC серверу
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gen.NewAuthServiceClient(conn)

	// Тестируем с невалидным токеном
	log.Println("Testing with invalid token...")
	resp, err := client.VerifyToken(context.Background(), &gen.VerifyTokenRequest{
		Token: "invalid-token",
	})
	if err != nil {
		log.Fatalf("gRPC call failed: %v", err)
	}
	log.Printf("Response: valid=%t, user_id=%d", resp.Valid, resp.UserId)

	// Тестируем с валидным токеном (замени на реальный JWT)
	log.Println("\nTesting with valid token...")
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ2ZXJzaW9uIjoxLCJpc3MiOiJhdXRoLXNlcnZpY2UiLCJleHAiOjE3NjE0ODYzNjksIm5iZiI6MTc2MTM5OTk2OSwiaWF0IjoxNzYxMzk5OTY5fQ.zptf2knVn0P9kNJiw_vKP_chGlNhEi3JzKHIHShLPJk" // Замени на реальный токен
	resp2, err := client.VerifyToken(context.Background(), &gen.VerifyTokenRequest{
		Token: validToken,
	})
	if err != nil {
		log.Fatalf("gRPC call failed: %v", err)
	}
	log.Printf("Response: valid=%t, user_id=%d", resp2.Valid, resp2.UserId)
}
