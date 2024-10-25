package dto

import "github.com/sanika-farm/sanika-farm-be/internal/domain/users/model"

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int    `json:"roleId"`
}

func (r *CreateUserRequest) ToModel() model.User {
	return model.User{
		Username: r.Username,
		Password: r.Password,
		RoleID:   r.RoleID,
	}
}
