package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"auth-service/internal/models"
	//"auth-service/internal/service"
)

type HTTPHandler struct {
	authService models.AuthService
}

func NewHTTPHandler(authService models.AuthService) *HTTPHandler {
	return &HTTPHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register data"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/register [post]
func (h *HTTPHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request data: " + err.Error(),
		})
		return
	}

	if err := h.authService.Register(req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{
		Message: "User registered successfully",
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login data"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/login [post]
func (h *HTTPHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request data: " + err.Error(),
		})
		return
	}

	response, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate user's JWT token
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.MessageResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/logout [post]
func (h *HTTPHandler) Logout(c *gin.Context) {
	token := h.extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Authorization token required",
		})
		return
	}

	if err := h.authService.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{
		Message: "Logged out successfully",
	})
}

// Me godoc
// @Summary Get current user info
// @Description Get information about currently authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.MeResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/me [get]
func (h *HTTPHandler) Me(c *gin.Context) {
	token := h.extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Authorization token required",
		})
		return
	}

	user, err := h.authService.GetMe(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.MeResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

// extractToken извлекает токен из заголовка Authorization
func (h *HTTPHandler) extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Формат: "Bearer {token}"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
