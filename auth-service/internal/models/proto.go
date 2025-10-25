package models

// TokenClaims - claims из JWT токена
type TokenClaims struct {
	UserID  int `json:"user_id"`
	Version int `json:"version"`
}
