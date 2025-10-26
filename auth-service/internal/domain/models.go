package domain

import "time"

type User struct {
	ID        int64
	Email     string
	PassHash  []byte
	CreatedAt time.Time
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
