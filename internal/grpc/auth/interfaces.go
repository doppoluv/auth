package auth

import "context"

type Auth interface {
	Login(ctx context.Context, username, password string) (token string, err error)
	Register(ctx context.Context, username, password, email string) (user_id int, err error)
	IsAdmin(ctx context.Context, user_id int) (is_admin bool, err error)
}
