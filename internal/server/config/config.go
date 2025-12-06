package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const dotEnvFile = ".env.server"

type Config struct {
	// Log level: debug, info, warn, error, fatal.
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	// gRPC server address.
	GRPCAddr string `env:"GRPC_ADDRESS,required,notEmpty"`
	// Database connection string.
	StorageDSN string `env:"DATABASE_URI,required,notEmpty"`
	// JWT secret key.
	JWTSecretKey string `env:"JWT_SECRET_KEY,required,notEmpty"`
	// JWT token TTL in seconds.
	JWTExpiration time.Duration `env:"JWT_EXPIRATION" envDefault:"24h"`
}

func LoadConfig() (*Config, error) {
	// Load .env.client file if exists (optional)
	_ = godotenv.Load(dotEnvFile)

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}
