package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sanika-farm/sanika-farm-be/internal/domain/users/model/dto"
	"github.com/sanika-farm/sanika-farm-be/transports/http/response"
)

// Register creates a new User.
// @Summary Create a new User.
// @Description This endpoint creates a new User.
// @Tags users
// @Param User body dto.CreateUserRequest true "The User to be created."
// @Produce json
// @Success 201 {object} response.Base{data=dto.CreateUserRequest}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/users/register [post]
func (h *UsersHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.UserService.CreateUser(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response.WithJSON(c, 201, req)
}
