package auth

import "fmt"

var (
	ErrInvalidCredentials    = fmt.Errorf("invalid credentials")
	ErrUsernameAlreadyExists = fmt.Errorf("username already exists")
	ErrEmailAlreadyExists    = fmt.Errorf("email already exists")
)
