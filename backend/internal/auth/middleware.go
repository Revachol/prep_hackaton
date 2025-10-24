package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Ключ для хранения userID в контексте.
type ctxKey string

const userIDKey ctxKey = "userID"

// UserIDFromContext — достаёт userID из контекста.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}

func AuthMiddleware(authManager AuthManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		tokenData, err := authManager.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("userID", tokenData.UserID)
		c.Next()
	}
}

// RequireAuth — обёртка для middleware в стиле "guard".
// Используется, если ты хочешь явно проверять авторизацию в хендлере.
//
// Пример:
//
//	func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
//	    userID, ok := auth.RequireAuth(r.Context())
//	    if !ok {
//	        http.Error(w, "unauthorized", http.StatusUnauthorized)
//	        return
//	    }
//	    fmt.Fprintf(w, "Hello, user %d", userID)
//	}
func RequireAuth(ctx context.Context) (int64, error) {
	id, ok := UserIDFromContext(ctx)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	return id, nil
}
