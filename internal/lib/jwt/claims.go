package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"auth/internal/domain/model"
)

type Claims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewClaims(user *model.User, duration time.Duration) Claims {
	expiresAt := time.Now().Add(duration)

	return Claims{
		UserId:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
}
