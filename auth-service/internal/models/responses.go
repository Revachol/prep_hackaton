package models

// LoginResponse represents successful login response
// @Description Login response with JWT token and user data
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// MeResponse represents current user data
// @Description Current user information
type MeResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// ErrorResponse represents error response
// @Description Error response with message
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse represents success message response
// @Description Success message response
type MessageResponse struct {
	Message string `json:"message"`
}
