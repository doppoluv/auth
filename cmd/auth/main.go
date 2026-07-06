package main

import (
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

}
