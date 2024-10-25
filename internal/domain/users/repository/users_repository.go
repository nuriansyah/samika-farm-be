package repository

import (
	"context"

	"github.com/sanika-farm/sanika-farm-be/internal/domain/users/model"
)

var (
	createUsers = struct {
		Query string
	}{
		Query: `INSERT INTO users (username, password, roleId) VALUES (?, ?, ?)`,
	}
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
}

func (r *UsersRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	tx, err := r.DB.Write.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(createUsers.Query, user.Username, user.Password, user.RoleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
