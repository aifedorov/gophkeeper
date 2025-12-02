package application

import (
	"context"

	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/container"
	"github.com/aifedorov/gophkeeper/internal/client/gui/root"
	"go.uber.org/zap"
)

type App struct {
	cfg      *config.Config
	logger   *zap.Logger
	services *container.Services
}

func NewApp(cfg *config.Config, logger *zap.Logger, services *container.Services) *App {
	return &App{
		services: services,
		cfg:      cfg,
		logger:   logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	gui := root.NewRoot(a.services)

	if err := gui.Run(); err != nil {
		a.logger.Error("failed to run gui", zap.Error(err))
		return err
	}

	return nil
}
