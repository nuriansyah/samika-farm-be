package services

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sanika-farm/sanika-farm-be/internal/domain/users/model/dto"
	"github.com/sanika-farm/sanika-farm-be/pkg/failure"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) error
}

func (s UsersServiceImpl) CreateUser(ctx context.Context, req *dto.CreateUserRequest) error {
	user := req.ToModel()
	err := s.UsersRepository.CreateUser(ctx, &user)
	if err != nil {
		if failure.GetCode(err) >= 500 {
			log.Error().Err(err).Msgf("Failed to create user: %v", err)
			return err
		}
		log.Warn().Err(err).Msgf("Failed to create user: %v", err)
		return err
	}
	return nil
}
