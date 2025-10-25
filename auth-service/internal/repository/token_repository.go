package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"auth-service/internal/models"
	"auth-service/pkg/redis"
)

type tokenRepository struct {
	redisClient *redis.Client
	expiration  time.Duration
}

func NewTokenRepository(redisClient *redis.Client, expiration time.Duration) models.TokenRepository {
	return &tokenRepository{
		redisClient: redisClient,
		expiration:  expiration,
	}
}

func (r *tokenRepository) getKey(userID int) string {
	return fmt.Sprintf("token_version:%d", userID)
}

func (r *tokenRepository) SetTokenVersion(ctx context.Context, userID int, version int) error {
	key := r.getKey(userID)
	return r.redisClient.Set(ctx, key, version, r.expiration)
}

func (r *tokenRepository) GetTokenVersion(ctx context.Context, userID int) (int, error) {
	key := r.getKey(userID)

	versionStr, err := r.redisClient.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, fmt.Errorf("invalid token version format: %w", err)
	}

	return version, nil
}

func (r *tokenRepository) DeleteTokenVersion(ctx context.Context, userID int) error {
	key := r.getKey(userID)
	return r.redisClient.Delete(ctx, key)
}

func (r *tokenRepository) IncrementTokenVersion(ctx context.Context, userID int) error {
	// Получаем текущую версию
	currentVersion, err := r.GetTokenVersion(ctx, userID)
	if err != nil {
		// Если записи нет, создаем новую с версией 1
		return r.SetTokenVersion(ctx, userID, 1)
	}

	// Увеличиваем версию
	return r.SetTokenVersion(ctx, userID, currentVersion+1)
}
