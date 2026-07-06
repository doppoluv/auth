package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	grpcauth "auth/internal/grpc/auth"
	"auth/internal/logger/interfaces"
)

type App struct {
	log  interfaces.Logger
	grpc *grpc.Server
	port int
}

func NewApp(
	log interfaces.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	grpcauth.Register(gRPCServer)

	return &App{
		log:  log,
		grpc: gRPCServer,
		port: port,
	}
}

func (a *App) Run() error {
	log := a.log

	log.Printf("Starting gRPC server on port %d", a.port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("listen tcp: %w", err)
	}

	log.Printf("gRPC server listening on %s", listener.Addr().String())

	if err := a.grpc.Serve(listener); err != nil {
		return fmt.Errorf("serve gRPC: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	log := a.log

	log.Printf("Stopping gRPC server")

	a.grpc.GracefulStop()

	log.Printf("gRPC server stopped")
}
