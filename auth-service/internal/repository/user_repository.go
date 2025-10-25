package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) models.UserRepository {
	return &userRepository{db: db}
}

// CreateUser создает нового пользователя
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (email, password, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`

	err = r.db.QueryRowContext(ctx, query,
		user.Email,
		string(hashedPassword),
		time.Now(),
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail возвращает пользователя по email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, created_at 
		FROM users 
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID возвращает пользователя по ID
func (r *userRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, password, created_at 
		FROM users 
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// VerifyPassword проверяет пароль пользователя
func (r *userRepository) VerifyPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
