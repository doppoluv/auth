package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"auth/internal/lib/jwt"
	"auth/internal/lib/logger"
	"auth/internal/storage"
)

var errInvalidCredentials = fmt.Errorf("invalid credentials")

type Auth struct {
	log          logger.Logger
	tokenTTL     time.Duration
	userSaver    UserSaver
	userProvider UserProvider
}

// NewAuth creates a new instance of the Auth service with the provided dependencies.
func NewAuth(
	log logger.Logger,
	tokenTTL time.Duration,
	userSaver UserSaver,
	userProvider UserProvider,
) *Auth {
	return &Auth{
		log:          log,
		tokenTTL:     tokenTTL,
		userSaver:    userSaver,
		userProvider: userProvider,
	}
}

// Login authenticates a user based on the provided username, password, and appID.
//
// It returns a JWT token if the authentication is successful.
func (a *Auth) Login(
	ctx context.Context,
	username, password string,
) (string, error) {
	log := a.log

	log.Printf("Authenticating user with username: %s", username)

	currentUser, err := a.userProvider.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", fmt.Errorf("get user by username: %w", errInvalidCredentials)
		}

		return "", fmt.Errorf("get user by username: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(currentUser.PasswordHash, []byte(password)); err != nil {
		return "", fmt.Errorf("compare password hash: %w", errInvalidCredentials)
	}

	claims := jwt.NewClaims(currentUser, a.tokenTTL)
	token, err := jwt.NewToken(claims)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	log.Printf("User %s authenticated successfully", username)

	return token, nil
}

// Register registers a new user with the provided email, username, and password.
//
// It returns the userId of the newly created user.
func (a *Auth) Register(
	ctx context.Context,
	email, username, password string,
) (int64, error) {
	log := a.log

	log.Printf("Registering user with username: %s", username)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("generate password hash: %w", err)
	}

	userId, err := a.userSaver.SaveUser(ctx, email, username, passwordHash)
	if err != nil {
		if errors.Is(err, storage.ErrUsernameAlreadyExists) {
			return 0, fmt.Errorf("save user: %w", storage.ErrUsernameAlreadyExists)
		}

		if errors.Is(err, storage.ErrEmailAlreadyExists) {
			return 0, fmt.Errorf("save user: %w", storage.ErrEmailAlreadyExists)
		}

		return 0, fmt.Errorf("save user: %w", err)
	}

	log.Printf("User %s registered successfully with userId: %d", username, userId)

	return userId, nil
}

// IsAdmin checks if a user with the given userId has administrative privileges.
func (a *Auth) IsAdmin(
	ctx context.Context,
	userId int64,
) (bool, error) {
	log := a.log

	log.Printf("Checking if user with userID: %d is an admin", userId)

	isAdmin, err := a.userProvider.IsUserAdmin(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("check if user is admin: %w", err)
	}

	return isAdmin, nil
}
