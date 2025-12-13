package container

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	client "github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/repository"
	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	grpcClient "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc"
	auth2 "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/auth"
	grpc "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/credential"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/storage"
)

type Services struct {
	AuthSrv       auth.Service
	CredsSrv      credential.Service
	TokenProvider client.TokenProvider

	grpcConn grpcClient.GRPCConnection
}

func NewServices(ctx context.Context, cfg *config.Config) (*Services, error) {
	sessionStore := storage.NewStorage()
	tokenProvider := auth.NewTokeProvider(sessionStore)
	conn, err := grpcClient.NewGRPCConnection(cfg.ServerAddr, tokenProvider)
	if err != nil {
		return nil, fmt.Errorf("container: failed to create grpc connection: %w", err)
	}

	authClient := auth2.NewAuthClient(conn.Conn())
	store := storage.NewStorage()
	repo := repository.NewRepository(ctx, store)
	authService := auth.NewService(authClient, repo)

	credClient := grpc.NewCredentialClient(conn.Conn())
	credService := credential.NewService(credClient)

	return &Services{
		AuthSrv:  authService,
		CredsSrv: credService,
		grpcConn: conn,
	}, nil
}

func (s *Services) Close() error {
	if s.grpcConn == nil {
		return fmt.Errorf("container: grpc connection is not initialized")
	}
	return s.grpcConn.Close()
}
