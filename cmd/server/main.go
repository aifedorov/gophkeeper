package main

import (
	"fmt"
	"os"

	"github.com/aifedorov/gophkeeper/internal/config"
	"github.com/aifedorov/gophkeeper/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = log.Sync()
	}()
}
