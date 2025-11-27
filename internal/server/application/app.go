package application

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/config"
	"github.com/aifedorov/gophkeeper/internal/server/domain/user"
	"github.com/aifedorov/gophkeeper/internal/server/domain/user/repository/db"
	"github.com/aifedorov/gophkeeper/pkg/posgres"
	"go.uber.org/zap"
)

type App struct {
	cfg    config.Config
	logger *zap.Logger
}

func NewApp(cfg config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run() error {
	ctx := context.Background()
	db := posgres.NewPosgresConnection(ctx, a.cfg.StorageDSN)
	err := db.Open()
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	userRepo := repository.NewRepository(ctx, db.DBPool(), a.logger)
	_ = user.NewService(userRepo, a.logger)

	return nil
}
