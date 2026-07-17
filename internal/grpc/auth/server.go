package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "auth/gen/go/auth/v1"
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
	if err := s.validateUsername(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate username: %v", err)
	}

	if err := s.validatePassword(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate password: %v", err)
	}

	if err := s.validateEmail(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate email: %v", err)
	}

	user_id, err := s.auth.Register(ctx, req.GetUsername(), req.GetPassword(), req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register: %v", err)
	}

	return &authv1.RegisterResponse{UserId: int64(user_id)}, nil
}

func (s *ServerAPI) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	if err := s.validateUsername(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate username: %v", err)
	}

	if err := s.validatePassword(); err != nil {
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
	if err := s.validateUserId(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validate user id: %v", err)
	}

	is_admin, err := s.auth.IsAdmin(ctx, int(req.GetUserId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "is_admin: %v", err)
	}

	return &authv1.IsAdminResponse{IsAdmin: is_admin}, nil
}

func (s *ServerAPI) validateEmail() error {
	// TODO: валидировать email
	return nil
}

func (s *ServerAPI) validateUsername() error {
	// TODO: валидировать username
	return nil
}

func (s *ServerAPI) validatePassword() error {
	// TODO: валидировать пароль
	return nil
}

func (s *ServerAPI) validateUserId() error {
	// TODO: валидировать user_id
	return nil
}
