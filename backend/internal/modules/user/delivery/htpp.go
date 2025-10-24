package delivery

import (
	"net/http"

	"github.com/Revachol/prep_hacakton/backend/internal/auth"
	"github.com/Revachol/prep_hacakton/backend/internal/modules/user/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *usecase.UserService
	authManager auth.AuthManager
}

// NewUserHandler — конструктор
func NewUserHandler(r *gin.Engine, userService *usecase.UserService, authManager auth.AuthManager) {
	h := &UserHandler{
		userService: userService,
		authManager: authManager,
	}

	// Публичные маршруты
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	// Группируем защищённые маршруты через middleware
	authGroup := r.Group("/", gin.WrapH(auth.Middleware(authManager, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))))
	authGroup.POST("/logout", h.Logout)
	authGroup.GET("/me", h.Me)
}
