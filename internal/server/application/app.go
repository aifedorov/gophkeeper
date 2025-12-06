package application

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/config"
	"github.com/aifedorov/gophkeeper/internal/server/domain/user"
	repository "github.com/aifedorov/gophkeeper/internal/server/domain/user/repository/db"
	server "github.com/aifedorov/gophkeeper/internal/server/infrastructure/grpc"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/jwt"
	"github.com/aifedorov/gophkeeper/pkg/posgres"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewApp(cfg *config.Config, logger *zap.Logger) *App {
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
	userSrv := user.NewService(userRepo, a.logger)
	jwtSrv := jwt.NewService(a.cfg.JWTSecretKey, a.cfg.JWTExpiration, a.logger)

	authServer := server.NewAuthServer(a.cfg, a.logger, userSrv, jwtSrv)
	grpcSrv := server.NewGRRPCServer(a.cfg, a.logger, grpc.NewServer(), authServer)

	if err := grpcSrv.Run(ctx); err != nil {
		return err
	}

	return nil
}
