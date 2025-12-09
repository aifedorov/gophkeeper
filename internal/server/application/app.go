package application

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/config"
	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	repository "github.com/aifedorov/gophkeeper/internal/server/domain/auth/repository/db"
	server "github.com/aifedorov/gophkeeper/internal/server/infrastructure/grpc"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/jwt"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/posgres"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
}

// NewApp creates a new instance of the application with the provided configuration and logger.
// It initializes the main application structure that orchestrates the server components.
func NewApp(cfg *config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

// Run initializes and starts the server application.
// It establishes database connection, initializes services (auth service, JWT service),
// creates the gRPC server with authentication handlers, and starts listening for requests.
// Returns an error if any initialization step fails.
func (a *App) Run() error {
	ctx := context.Background()
	db := posgres.NewPosgresConnection(ctx, a.cfg.StorageDSN)
	err := db.Open()
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	userRepo := repository.NewRepository(db.DBPool(), a.logger)
	userSrv := auth.NewService(userRepo, a.logger)
	jwtSrv := jwt.NewService(a.cfg.JWTSecretKey, a.cfg.JWTExpiration, a.logger)

	authServer := server.NewAuthServer(a.cfg, a.logger, userSrv, jwtSrv)
	grpcSrv := server.NewGRRPCServer(a.cfg, a.logger, grpc.NewServer(), authServer)

	if err := grpcSrv.Run(ctx); err != nil {
		return err
	}

	return nil
}
