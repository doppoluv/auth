package main

import (
	"auth/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Configuration loaded: %+v\n", cfg)
}
