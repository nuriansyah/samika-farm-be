package repository

import "github.com/sanika-farm/sanika-farm-be/infras"

// UsersRepository is the interface for repository.
type UsersRepository interface {
	UserRepository
}

type UsersRepositoryImpl struct {
	DB *infras.PostgresConn
}

func ProvideUsersRepository(db *infras.PostgresConn) *UsersRepositoryImpl {
	return &UsersRepositoryImpl{
		DB: db,
	}
}
