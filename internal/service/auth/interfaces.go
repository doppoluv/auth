package auth

import (
	"context"

	"auth/internal/domain/model"
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email, username string,
		passwordHash []byte,
	) (int64, error)
}

type UserProvider interface {
	GetUserByUsername(
		ctx context.Context,
		username string,
	) (*model.User, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*model.User, error)
	IsUserAdmin(
		ctx context.Context,
		userId int64,
	) (bool, error)
}
