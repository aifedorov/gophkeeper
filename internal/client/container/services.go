package container

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	grpcClient "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/memory"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Services struct {
	AuthSrv auth.Service

	grpcConn *grpc.ClientConn
	logger   *zap.Logger
}

func NewServices(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*Services, error) {
	conn, err := createGRPCConnection(ctx, cfg.ServerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc connection: %w", err)
	}

	authClient := grpcClient.NewAuthClient(conn)
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
	if s.grpcConn != nil {
		return s.grpcConn.Close()
	}
	return nil
}

func createGRPCConnection(ctx context.Context, serverAddr string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		// TODO: Add TLS
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	return conn, nil
}
