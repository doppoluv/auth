package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	"auth/internal/app"
	"auth/internal/config"
	"auth/internal/lib/logger"
)

func main() {
	log := logger.NewLogger()

	err := godotenv.Load()
	if err != nil {
		log.Warningf("error load .env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load configuration: %w", err)
	}

	log.Infof("Configuration loaded\n")

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

	log.Infof("Application stopped")
}
