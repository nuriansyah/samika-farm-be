package services

import (
	"github.com/sanika-farm/sanika-farm-be/configs"
	"github.com/sanika-farm/sanika-farm-be/internal/domain/users/repository"
)

type UsersService interface {
	UserService
}

type UsersServiceImpl struct {
	UsersRepository repository.UsersRepository
	cfg             *configs.Config
}

func ProvideUsersService(repo repository.UsersRepository, cfg *configs.Config) *UsersServiceImpl {
	return &UsersServiceImpl{
		UsersRepository: repo,
		cfg:             cfg,
	}
}
