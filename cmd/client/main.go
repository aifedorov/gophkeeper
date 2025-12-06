package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aifedorov/gophkeeper/internal/client/application"
	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/container"
	"github.com/aifedorov/gophkeeper/internal/client/version"
	"github.com/aifedorov/gophkeeper/pkg/logger"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version.Info())
		os.Exit(0)
	}

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

	app := application.NewApp(cfg, log, services)
	if err := app.RunCLI(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to run application: %v\n", err)
		os.Exit(1)
	}
}
