package auth

import (
	"context"

	"auth/internal/domain/model"
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		username string,
		passwordHash []byte,
	) (userID int64, err error)
}

type UserProvider interface {
	GetUserByUsername(
		ctx context.Context,
		username string,
	) (user model.User, err error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (user model.User, err error)
	IsUserAdmin(
		ctx context.Context,
		userID int64,
	) (isAdmin bool, err error)
}

type AppProvider interface {
	GetAppByID(
		ctx context.Context,
		appID int64,
	) (app model.App, err error)
}
