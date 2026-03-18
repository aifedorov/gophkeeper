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
	// gRPC binary address.
	GRPCAddr string `env:"GRPC_ADDRESS,required,notEmpty"`
	// Database connection string.
	StorageDSN string `env:"DATABASE_URI,required,notEmpty"`
	// JWT secret key.
	JWTSecretKey string `env:"JWT_SECRET_KEY,required,notEmpty"`
	// JWT token TTL in seconds.
	JWTExpiration time.Duration `env:"JWT_EXPIRATION" envDefault:"24h"`
	// TLS certificate file path.
	TLSCertPath string `env:"TLS_CERT_PATH" envDefault:"certs/server-cert.pem"`
	// TLS private key file path.
	TLSKeyPath string `env:"TLS_KEY_PATH" envDefault:"certs/server-key.pem"`
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
