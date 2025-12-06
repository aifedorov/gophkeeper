package application

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/cli/register"
	cliroot "github.com/aifedorov/gophkeeper/internal/client/cli/root"
	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/container"
	guiroot "github.com/aifedorov/gophkeeper/internal/client/gui/root"
	"go.uber.org/zap"
)

type App struct {
	cfg      *config.Config
	logger   *zap.Logger
	services *container.Services
	rootCmd  *cliroot.RootCommand
}

func NewApp(cfg *config.Config, logger *zap.Logger, services *container.Services) *App {
	rootCmd := cliroot.NewCommand()
	registerCmd := register.NewCommand(services.AuthSrv)
	rootCmd.AddCommand(registerCmd.GetCommand())

	return &App{
		services: services,
		cfg:      cfg,
		logger:   logger,
		rootCmd:  rootCmd,
	}
}

func (a *App) RunCLI(ctx context.Context) error {
	if err := a.rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to run cli: %w", err)
	}
	return nil
}

func (a *App) RunGUI(ctx context.Context) error {
	gui := guiroot.NewRoot(a.services)

	if err := gui.Run(); err != nil {
		a.logger.Error("failed to run gui", zap.Error(err))
		return err
	}

	return nil
}
