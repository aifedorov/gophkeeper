package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const dotEnvFile = ".env"

type Config struct {
	// Log level: debug, info, warn, error, fatal.
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	// gRPC server address.
	RunAddr string `env:"GRPC_ADDRESS" envDefault:"localhost:9090"`
	// Database connection string.
	StorageDSN string `env:"DATABASE_URI,required,notEmpty"`
}

func LoadConfig() (Config, error) {
	err := godotenv.Load(dotEnvFile)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
