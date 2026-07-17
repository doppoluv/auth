package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"k8s.io/utils/env"
)

func NewToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(env.GetString("AUTH_JWT_SECRET", ""))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}
