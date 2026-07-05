package main

import (
	"log"

	"auth/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	log.Printf("Configuration loaded: %+v\n", cfg) // remove this line in prod

}
