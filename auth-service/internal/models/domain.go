package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenVersion struct {
	UserID    int       `json:"user_id"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}
