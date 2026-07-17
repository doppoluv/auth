package auth

import "context"

type Auth interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, email, username, password string) (int64, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}
