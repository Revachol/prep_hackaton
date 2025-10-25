package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"auth-service/internal/models"
	//"auth-service/internal/service"
)

type AuthMiddleware struct {
	authService models.AuthService
}

func NewAuthMiddleware(authService models.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// AuthRequired middleware проверяет JWT токен
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization token required",
			})
			c.Abort()
			return
		}

		claims, err := m.authService.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Сохраняем user_id в контекст для использования в обработчиках
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// extractToken извлекает токен из заголовка Authorization
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
