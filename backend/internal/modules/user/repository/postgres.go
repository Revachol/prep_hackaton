package repository

import (
	"context"
	"time"

	"github.com/Revachol/prep_hackaton/backend/internal/modules/user/entity"
	"github.com/jmoiron/sqlx"
)

// UserRepository — интерфейс для работы с пользователями.
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	UpdateTokenVersion(ctx context.Context, id int64, version int) error
}

// PostgresUserRepository — реализация через PostgreSQL (sqlx).
type PostgresUserRepository struct {
	db *sqlx.DB
}

// NewPostgresUserRepository — конструктор.
func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// ============================
// CreateUser
// ============================
func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
	INSERT INTO users (email, password_hash, token_version, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.TokenVersion,
		time.Now(),
		time.Now(),
	).Scan(&user.ID)
	return err
}

// ============================
// GetByEmail
// ============================
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	query := `SELECT id, email, password_hash, token_version, created_at, updated_at FROM users WHERE email=$1`
	err := r.db.GetContext(ctx, user, query, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ============================
// GetByID
// ============================
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	user := &entity.User{}
	query := `SELECT id, email, password_hash, token_version, created_at, updated_at FROM users WHERE id=$1`
	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ============================
// UpdateTokenVersion
// ============================
func (r *PostgresUserRepository) UpdateTokenVersion(ctx context.Context, id int64, version int) error {
	query := `UPDATE users SET token_version=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.ExecContext(ctx, query, version, time.Now(), id)
	return err
}
