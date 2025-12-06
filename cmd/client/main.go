package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aifedorov/gophkeeper/internal/client/application"
	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/container"
	"github.com/aifedorov/gophkeeper/pkg/logger"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize config: %v\n", err)
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

	services, err := container.NewServices(ctx, cfg, log)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create services: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := services.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to close services: %v\n", err)
		}
	}()

	app, err := application.NewApp(cfg, log, services)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create application: %v\n", err)
		os.Exit(1)
	}
	if err := app.RunCLI(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to run application: %v\n", err)
		os.Exit(1)
	}
}
