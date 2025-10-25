package delivery

import (
	"net/http"

	"backend/internal/auth"
	"backend/internal/modules/user/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *usecase.UserService
}

// NewUserHandler подключает маршруты Gin
func NewUserHandler(r *gin.Engine, userService *usecase.UserService, authManager auth.AuthManager) {
	h := &UserHandler{
		userService: userService,
	}

	// Публичные маршруты
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	// Защищённые маршруты через middleware
	authGroup := r.Group("/", auth.AuthMiddleware(authManager))
	authGroup.POST("/logout", h.Logout)
	authGroup.GET("/me", h.Me)
}

// ===================== Register =====================
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.userService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ===================== Login =====================
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ===================== Logout =====================
func (h *UserHandler) Logout(c *gin.Context) {
	userID, exists := auth.UserIDFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.userService.Logout(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// ===================== Me =====================
func (h *UserHandler) Me(c *gin.Context) {
	userID, exists := auth.UserIDFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.Me(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}
