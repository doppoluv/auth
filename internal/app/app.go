package app

import (
	appgrpc "auth/internal/app/grpc"
	"auth/internal/logger/interfaces"
)

type App struct {
	GRPCServer *appgrpc.App
}

func NewApp(
	log interfaces.Logger,
	port int,
	storagePath string,
	tokenTTL string,
) *App {
	// TODO: инициализация хранилища
	// TODO: инициализация AUTH сервиса

	GRPCServer := appgrpc.NewApp(log, port)
	return &App{
		GRPCServer: GRPCServer,
	}
}
