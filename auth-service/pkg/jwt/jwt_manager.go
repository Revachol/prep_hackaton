package jwt

import (
	"errors"
	"fmt"
	"time"

	"auth-service/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// Claims - кастомные claims для нашего JWT токена
type Claims struct {
	UserID  int `json:"user_id"`
	Version int `json:"version"`
	jwt.RegisteredClaims
}

// GenerateToken создает новый JWT токен для пользователя
func (m *JWTManager) GenerateToken(userID int, version int) (string, error) {
	// Создаем claims
	claims := &Claims{
		UserID:  userID,
		Version: version,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	// Создаем токен с claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken проверяет JWT токен и возвращает claims
func (m *JWTManager) VerifyToken(tokenString string) (*models.TokenClaims, error) {
	// Парсим токен с claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Извлекаем claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Проверяем expiration
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token has expired")
	}

	return &models.TokenClaims{
		UserID:  claims.UserID,
		Version: claims.Version,
	}, nil
}

// ExtractClaims извлекает claims из токена без полной валидации (для отладки)
func (m *JWTManager) ExtractClaims(tokenString string) (*Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
