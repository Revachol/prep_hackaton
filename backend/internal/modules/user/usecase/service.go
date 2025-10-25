package usecase

import (
	"context"
	"errors"

	"backend/internal/auth"
	"backend/internal/modules/user/entity"
	"backend/internal/modules/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo        repository.UserRepository
	authManager auth.AuthManager
}

// NewUserService — конструктор
func NewUserService(repo repository.UserRepository, authManager auth.AuthManager) *UserService {
	return &UserService{
		repo:        repo,
		authManager: authManager,
	}
}

// ===================== Register =====================
func (s *UserService) Register(ctx context.Context, email, password string) (string, error) {
	// Проверяем, есть ли пользователь
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return "", errors.New("user already exists")
	}

	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &entity.User{
		Email:        email,
		PasswordHash: string(hash),
		TokenVersion: 1,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", err
	}

	// Генерируем JWT
	token, err := s.authManager.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ===================== Login =====================
func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Генерируем JWT
	token, err := s.authManager.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ===================== Logout =====================
func (s *UserService) Logout(ctx context.Context, userID int64) error {
	// Получаем пользователя
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Инвалидируем токены — увеличиваем token_version
	newVersion := user.TokenVersion + 1
	if err := s.repo.UpdateTokenVersion(ctx, userID, newVersion); err != nil {
		return err
	}

	// Также можно очистить Redis (если используем кэш версий)
	if s.authManager.Invalidate != nil {
		if err := s.authManager.Invalidate(userID); err != nil {
			return err
		}
	}

	return nil
}

// ===================== Me =====================
func (s *UserService) Me(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Не возвращаем пароль и token_version
	return &entity.User{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
