package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"auth-service/internal/models"
	"auth-service/pkg/jwt"
)

type authService struct {
	jwtManager *jwt.JWTManager
	tokenRepo  models.TokenRepository
	userRepo   models.UserRepository // Правильный тип!
}

// ФИКС: меняем models.TokenRepository на models.UserRepository для userRepo
func NewAuthService(jwtManager *jwt.JWTManager, tokenRepo models.TokenRepository, userRepo models.UserRepository) models.AuthService {
	return &authService{
		jwtManager: jwtManager,
		tokenRepo:  tokenRepo,
		userRepo:   userRepo, // Теперь правильный тип
	}
}

func (s *authService) Register(email, password string) error {
	log.Printf("Register called: email=%s", email)

	ctx := context.Background()

	// Проверяем не существует ли уже пользователь
	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return errors.New("user already exists")
	}

	// Создаем нового пользователя
	user := &models.User{
		Email:    email,
		Password: password,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("🟢 User registered successfully: id=%d, email=%s", user.ID, user.Email)
	return nil
}

func (s *authService) Login(email, password string) (*models.LoginResponse, error) {
	log.Printf("Login called: email=%s", email)

	ctx := context.Background()

	// Ищем пользователя в БД
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Проверяем пароль
	if err := s.userRepo.VerifyPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Получаем или создаем версию токена
	version, err := s.tokenRepo.GetTokenVersion(ctx, user.ID)
	if err != nil {
		log.Printf("🔴 No token version found, creating new one: %v", err)
		version = 1
		if err := s.tokenRepo.SetTokenVersion(ctx, user.ID, version); err != nil {
			return nil, fmt.Errorf("failed to set token version: %w", err)
		}
		log.Printf("🟢 Created new token version: %d for user %d", version, user.ID)
	} else {
		log.Printf("🟢 Found existing token version: %d for user %d", version, user.ID)
	}

	// Генерируем JWT токен
	token, err := s.jwtManager.GenerateToken(user.ID, version)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("🟢 Generated JWT token for user %d with version %d", user.ID, version)

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *authService) Logout(tokenString string) error {
	log.Printf("Logout called")

	// Верифицируем токен чтобы получить user_id
	claims, err := s.jwtManager.VerifyToken(tokenString)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// Создаем контекст
	ctx := context.Background()

	// Инкрементируем версию токена (инвалидируем все старые токены)
	err = s.tokenRepo.IncrementTokenVersion(ctx, claims.UserID)
	if err != nil {
		return fmt.Errorf("failed to invalidate token: %w", err)
	}

	log.Printf("Logout successful for user_id=%d", claims.UserID)
	return nil
}

func (s *authService) VerifyToken(tokenString string) (*models.TokenClaims, error) {
	log.Printf("VerifyToken called")

	// Верифицируем подпись и expiration токена
	claims, err := s.jwtManager.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Создаем контекст
	ctx := context.Background()

	// Проверяем версию токена в Redis
	currentVersion, err := s.tokenRepo.GetTokenVersion(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("token version not found")
	}

	// Если версия в токене не совпадает с текущей - токен невалиден
	if claims.Version != currentVersion {
		return nil, errors.New("token has been invalidated")
	}

	return claims, nil
}

func (s *authService) GetMe(tokenString string) (*models.User, error) {
	log.Printf("GetMe called")

	// Верифицируем токен
	claims, err := s.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// ФИКС: Получаем реального пользователя из базы данных
	ctx := context.Background()
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Важно: не возвращаем хэшированный пароль
	user.Password = ""

	log.Printf("🟢 GetMe: returned user id=%d, email=%s", user.ID, user.Email)
	return user, nil
}
