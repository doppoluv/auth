package auth

import (
	"context"

	"google.golang.org/grpc"

	authv1 "auth/gen/go/auth/v1"
)

type ServerAPI struct {
	authv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &ServerAPI{})
}

func (s *ServerAPI) Register(
	ctx context.Context,
	req *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	// Implementation for Register method
	return &authv1.RegisterResponse{}, nil
}

func (s *ServerAPI) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	// Implementation for Login method
	return &authv1.LoginResponse{}, nil
}

func (s *ServerAPI) IsAdmin(
	ctx context.Context,
	req *authv1.IsAdminRequest,
) (*authv1.IsAdminResponse, error) {
	// Implementation for IsAdmin method
	return &authv1.IsAdminResponse{}, nil
}
