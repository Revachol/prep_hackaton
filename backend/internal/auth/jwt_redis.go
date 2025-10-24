package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// SessionJWTManager — реализация AuthManager, использующая JWT + Redis.
type SessionJWTManager struct {
	secret []byte
	ttl    time.Duration
	redis  *redis.Client
}

// NewSessionJWTManager — конструктор.
func NewSessionJWTManager(secret string, ttl time.Duration, redis *redis.Client) *SessionJWTManager {
	return &SessionJWTManager{
		secret: []byte(secret),
		ttl:    ttl,
		redis:  redis,
	}
}

// =============================
// 1. Генерация токена
// =============================
func (m *SessionJWTManager) GenerateToken(userID int64) (string, error) {
	ctx := context.Background()

	// Создаём новую версию токена — timestamp
	version := time.Now().UnixNano()

	// Сохраняем версию в Redis
	sessionKey := m.sessionKey(userID)
	if err := m.redis.Set(ctx, sessionKey, version, m.ttl).Err(); err != nil {
		return "", fmt.Errorf("redis set error: %w", err)
	}

	// Формируем JWT
	claims := jwt.MapClaims{
		"user_id": userID,
		"ver":     version,
		"exp":     time.Now().Add(m.ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenStr, nil
}

// =============================
// 2. Проверка токена
// =============================
func (m *SessionJWTManager) ParseToken(tokenStr string) (*TokenData, error) {
	ctx := context.Background()

	// Парсим токен
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	// Проверка валидности
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id claim")
	}
	version, ok := claims["ver"].(float64)
	if !ok {
		return nil, errors.New("invalid version claim")
	}

	// Проверяем актуальную версию в Redis
	sessionKey := m.sessionKey(int64(userID))
	currentVersion, err := m.redis.Get(ctx, sessionKey).Int64()
	if err == redis.Nil {
		return nil, errors.New("token invalidated or expired")
	} else if err != nil {
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	if currentVersion != int64(version) {
		return nil, errors.New("token invalidated")
	}

	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	return &TokenData{
		UserID: int64(userID),
		Exp:    exp,
	}, nil
}

// =============================
// 3. Инвалидация токена
// =============================
func (m *SessionJWTManager) Invalidate(userID int64) error {
	ctx := context.Background()
	sessionKey := m.sessionKey(userID)
	return m.redis.Del(ctx, sessionKey).Err()
}

// =============================
// Вспомогательные функции
// =============================

func (m *SessionJWTManager) sessionKey(userID int64) string {
	return fmt.Sprintf("user:%d:token_version", userID)
}
