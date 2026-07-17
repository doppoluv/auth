package app

import (
	"time"

	appgrpc "auth/internal/app/grpc"
	"auth/internal/lib/logger"
	"auth/internal/service/auth"
	"auth/internal/storage/sqlite"
)

type App struct {
	GRPCServer *appgrpc.App
}

func NewApp(
	log logger.Logger,
	port int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.NewStorage(storagePath)
	if err != nil {
		panic(err) // TODO: убрать панику
	}

	authService := auth.NewAuth(log, tokenTTL, storage, storage)

	GRPCServer := appgrpc.NewApp(port, log, authService)
	return &App{
		GRPCServer: GRPCServer,
	}
}
