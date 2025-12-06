package container

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	grpcClient "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/memory"
	"go.uber.org/zap"
)

type Services struct {
	AuthSrv auth.Service

	grpcConn grpcClient.GRPCConnection
	logger   *zap.Logger
}

func NewServices(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*Services, error) {
	conn, err := grpcClient.NewGRPCConnection(cfg.ServerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc connection: %w", err)
	}

	authClient := grpcClient.NewAuthClient(conn.Conn())
	storage := memory.NewStorage()
	repo := auth.NewRepository(ctx, logger, storage)
	authService := auth.NewService(authClient, repo)

	return &Services{
		AuthSrv:  authService,
		grpcConn: conn,
		logger:   logger,
	}, nil
}

func (s *Services) Close() error {
	return s.grpcConn.Close()
}
