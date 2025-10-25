package models

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	VerifyPassword(hashedPassword, plainPassword string) error // Добавляем этот метод
}

type TokenRepository interface {
	SetTokenVersion(ctx context.Context, userID int, version int) error
	GetTokenVersion(ctx context.Context, userID int) (int, error)
	DeleteTokenVersion(ctx context.Context, userID int) error
	IncrementTokenVersion(ctx context.Context, userID int) error
}

type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (*LoginResponse, error)
	Logout(token string) error
	VerifyToken(token string) (*TokenClaims, error)
	GetMe(token string) (*User, error)
}
