package application

import (
	"context"
	"fmt"

	cliroot "github.com/aifedorov/gophkeeper/internal/client/cli/root"
	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/container"
)

type App struct {
	cfg      *config.Config
	services *container.Services
	rootCmd  *cliroot.Command
}

func NewApp(cfg *config.Config, services *container.Services) (*App, error) {
	rootCmd, err := cliroot.NewCommand(services.AuthSrv, services.CredsSrv, services.BinarySrv, services.CardSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create root command: %w", err)
	}

	return &App{
		services: services,
		cfg:      cfg,
		rootCmd:  rootCmd,
	}, nil
}

func (a *App) RunCLI(ctx context.Context) error {
	if err := a.rootCmd.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("failed to run cli: %w", err)
	}
	return nil
}
