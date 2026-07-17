package main

import (
	"os"
	"os/signal"
	"time"

	"auth/internal/app"
	"auth/internal/config"
	"auth/internal/lib/logger"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load configuration: %w", err)
	}

	log.Printf("Configuration loaded\n")

	duration, err := time.ParseDuration(cfg.TokenTTL)
	if err != nil {
		log.Fatalf("parse duration: %w", err)
	}

	application := app.NewApp(log, cfg.GRPC.Port, cfg.StoragePath, duration)

	go func() {
		if err := application.GRPCServer.Run(); err != nil {
			log.Fatalf("run gRPC server: %w", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	application.GRPCServer.Stop()

	log.Printf("Application stopped")
}
