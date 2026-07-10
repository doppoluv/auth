package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"auth/internal/domain/model"
)

func NewToken(claims Claims, app model.App) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}
