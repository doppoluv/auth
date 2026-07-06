package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"auth/internal/logger"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	log.Printf("Configuration loaded: %+v\n", cfg) // remove this line in prod

	application := app.NewApp(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	if err := application.GRPCServer.Run(); err != nil {
		log.Fatalf("run gRPC server: %v", err)
	}
}
