package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sanika-farm/sanika-farm-be/internal/domain/users/services"
)

// UsersHandler is the HTTP handler for Users domain.
type UsersHandler struct {
	UserService services.UsersService
}

// ProvideUsersHandler is the provider for this handler.
func ProvideUsersHandler(svcUser services.UsersService) UsersHandler {
	return UsersHandler{
		UserService: svcUser,
	}
}

func (h *UsersHandler) Router(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("/register", h.CreateUser)
	}
}
