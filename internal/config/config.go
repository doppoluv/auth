package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" validate:"required,oneof=local production"`
	StoragePath string        `yaml:"storage_path" validate:"required"`
	TokenTTL    time.Duration `yaml:"token_ttl" validate:"required"`
	GRPC        GRPCConfig    `yaml:"grpc" validate:"required"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func Load() (*Config, error) {
	configPath := fetchConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("check configuration file: %w", err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}

	return &cfg, nil
}

func LoadByPath(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("check configuration file: %w", err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}

	return &cfg, nil
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to the configuration file")
	flag.Parse()

	return path
}
