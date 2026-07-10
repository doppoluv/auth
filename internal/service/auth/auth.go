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

var (
	errInvalidCredentials = fmt.Errorf("invalid credentials")
	errInvalidAppID       = fmt.Errorf("invalid app ID")
)

type Auth struct {
	log          logger.Logger
	tokenTTL     time.Duration
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
}

// NewAuth creates a new instance of the Auth service with the provided dependencies.
func NewAuth(
	log logger.Logger,
	tokenTTL time.Duration,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
) *Auth {
	return &Auth{
		log:          log,
		tokenTTL:     tokenTTL,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
	}
}

// LoginUser authenticates a user based on the provided username, password, and appID.
//
// It returns a JWT token if the authentication is successful.
func (a *Auth) LoginUser(
	ctx context.Context,
	username, password string,
	appID int64,
) (string, error) {
	log := a.log

	log.Printf("Authenticating user with username: %s for appID: %d", username, appID)

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

	currentApp, err := a.appProvider.GetAppByID(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("get app by ID: %w", err)
	}

	log.Printf("User %s authenticated successfully for app: %s", username, currentApp.Name)

	claims := jwt.NewClaims(currentUser, currentApp, a.tokenTTL)
	token, err := jwt.NewToken(claims, currentApp)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

// RegisterUser registers a new user with the provided email, username, and password.
//
// It returns the userID of the newly created user.
func (a *Auth) RegisterUser(
	ctx context.Context,
	email, username, password string,
) (int64, error) {
	log := a.log

	log.Printf("Registering user with username: %s", username)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("generate password hash: %w", err)
	}

	userID, err := a.userSaver.SaveUser(ctx, email, username, passwordHash)
	if err != nil {
		if errors.Is(err, storage.ErrUsernameAlreadyExists) {
			return 0, fmt.Errorf("save user: %w", storage.ErrUsernameAlreadyExists)
		}

		if errors.Is(err, storage.ErrEmailAlreadyExists) {
			return 0, fmt.Errorf("save user: %w", storage.ErrEmailAlreadyExists)
		}

		return 0, fmt.Errorf("save user: %w", err)
	}

	log.Printf("User %s registered successfully with userID: %d", username, userID)

	return userID, nil
}

// IsUserAdmin checks if a user with the given userID has administrative privileges.
func (a *Auth) IsUserAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	log := a.log

	log.Printf("Checking if user with userID: %d is an admin", userID)

	isAdmin, err := a.userProvider.IsUserAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return false, fmt.Errorf("check if user is admin: %w", errInvalidAppID)
		}

		return false, fmt.Errorf("check if user is admin: %w", err)
	}

	return isAdmin, nil
}
