package storage

import "fmt"

var (
	ErrUsernameAlreadyExists = fmt.Errorf("username already exists")
	ErrEmailAlreadyExists    = fmt.Errorf("email already exists")
	ErrUserNotFound          = fmt.Errorf("user not found")
	ErrAppNotFound           = fmt.Errorf("app not found")
)
