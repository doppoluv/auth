package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"auth/internal/domain/model"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	AppID    int64  `json:"app_id"`
	jwt.RegisteredClaims
}

func NewClaims(user model.User, app model.App, duration time.Duration) Claims {
	expiresAt := time.Now().Add(duration)

	return Claims{
		UserID:   user.ID,
		Username: user.Username,
		AppID:    app.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
}
