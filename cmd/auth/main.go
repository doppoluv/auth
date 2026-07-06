package main

import (
	"os"
	"os/signal"

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

	log.Printf("Configuration loaded\n")

	application := app.NewApp(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go func() {
		if err := application.GRPCServer.Run(); err != nil {
			log.Fatalf("run gRPC server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	application.GRPCServer.Stop()

	log.Printf("Application stopped")
}
