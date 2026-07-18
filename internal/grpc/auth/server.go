package auth

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "auth/gen/go/auth/v1"
)

var (
	errInvalidEmail    = fmt.Errorf("invalid email")
	errInvalidUserId   = fmt.Errorf("invalid user id")
	errInvalidUsername = fmt.Errorf("invalid username")
	errInvalidPassword = fmt.Errorf("invalid password")

	// username: only Latin letters, digits, underscore and hyphen are allowed,
	// first char must be alphanumeric, length is 3..32 characters.
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{2,31}$`)

	// password: only Latin letters and digits are allowed, and the total
	// length must be at least 8 characters.
	passwordPattern = regexp.MustCompile(`^[A-Za-z\d]{8,}$`)
)

type ServerAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Register(
	ctx context.Context,
	req *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	if err := s.validateUsername(req.GetUsername()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate username: %v", err)
	}

	if err := s.validatePassword(req.GetPassword()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate password: %v", err)
	}

	if err := s.validateEmail(req.GetEmail()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate email: %v", err)
	}

	user_id, err := s.auth.Register(ctx, req.GetEmail(), req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register: %v", err)
	}

	return &authv1.RegisterResponse{UserId: int64(user_id)}, nil
}

func (s *ServerAPI) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	if err := s.validateUsername(req.GetUsername()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate username: %v", err)
	}

	if err := s.validatePassword(req.GetPassword()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate password: %v", err)
	}

	token, err := s.auth.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login: %v", err)
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *ServerAPI) IsAdmin(
	ctx context.Context,
	req *authv1.IsAdminRequest,
) (*authv1.IsAdminResponse, error) {
	if err := s.validateUserId(req.GetUserId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate user id: %v", err)
	}

	is_admin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "is_admin: %v", err)
	}

	return &authv1.IsAdminResponse{IsAdmin: is_admin}, nil
}

func (s *ServerAPI) validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errInvalidEmail
	}

	return nil
}

func (s *ServerAPI) validateUsername(username string) error {
	if strings.TrimSpace(username) == "" {
		return errInvalidUsername
	}

	if !usernamePattern.MatchString(username) {
		return errInvalidUsername
	}

	return nil
}

func (s *ServerAPI) validatePassword(password string) error {
	if !passwordPattern.MatchString(password) {
		return errInvalidPassword
	}

	var hasUpper, hasLower, hasDigit bool
	for _, ch := range password {
		switch {
		case 'A' <= ch && ch <= 'Z':
			hasUpper = true
		case 'a' <= ch && ch <= 'z':
			hasLower = true
		case '0' <= ch && ch <= '9':
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return errInvalidPassword
	}

	return nil
}

func (s *ServerAPI) validateUserId(userId int64) error {
	if userId < 0 {
		return errInvalidUserId
	}

	return nil
}
