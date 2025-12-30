package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const dotEnvFile = ".env.client"

type Config struct {
	// Log level: debug, info, warn, error, fatal.
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	// gRPC binary address.
	ServerAddr string `env:"SERVER_ADDRESS,required,notEmpty"`
	// File storage root path.
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"storage/files/"`
}

func LoadConfig() (*Config, error) {
	// Load .env.client storage if exists (optional)
	_ = godotenv.Load(dotEnvFile)

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}
