package handlers

import (
	"context"
	"log"

	"auth-service/gen"
	"auth-service/internal/models"
)

type GRPCHandler struct {
	gen.UnimplementedAuthServiceServer
	authService models.AuthService
}

func NewGRPCHandler(authService models.AuthService) *GRPCHandler {
	return &GRPCHandler{
		authService: authService,
	}
}

// VerifyToken реализует gRPC метод для проверки токенов
func (h *GRPCHandler) VerifyToken(ctx context.Context, req *gen.VerifyTokenRequest) (*gen.VerifyTokenResponse, error) {
	log.Printf("gRPC VerifyToken called")

	// Используем наш существующий сервис для проверки токена
	claims, err := h.authService.VerifyToken(req.Token)
	if err != nil {
		log.Printf("gRPC VerifyToken: invalid token - %v", err)
		return &gen.VerifyTokenResponse{
			Valid:  false,
			UserId: 0,
		}, nil
	}

	log.Printf("gRPC VerifyToken: valid token for user_id=%d", claims.UserID)
	return &gen.VerifyTokenResponse{
		Valid:  true,
		UserId: int64(claims.UserID),
	}, nil
}
