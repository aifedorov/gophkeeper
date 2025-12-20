package container

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	authinterfaces "github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/repository"
	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	grpcClient "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/client"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/storage"
	"github.com/aifedorov/gophkeeper/pkg/filestorage"
	"go.uber.org/zap"
)

type Services struct {
	AuthSrv       auth.Service
	CredsSrv      credential.Service
	BinarySrv     binary.Service
	TokenProvider authinterfaces.SessionProvider
	grpcConn      grpcClient.GRPCConnection
}

func NewServices(ctx context.Context, cfg *config.Config) (*Services, error) {
	sessionStore := storage.NewStorage()
	fileStore := filestorage.NewFileStorage(zap.NewNop())
	tokenProvider := auth.NewSessionProvider(sessionStore)
	conn, err := grpcClient.NewGRPCConnection(cfg.ServerAddr, tokenProvider)
	if err != nil {
		return nil, fmt.Errorf("container: failed to create grpc connection: %w", err)
	}

	authClient := client.NewAuthClient(conn.Conn())
	store := storage.NewStorage()
	repo := repository.NewRepository(ctx, store)
	authService := auth.NewService(authClient, repo)

	credClient := client.NewCredentialClient(conn.Conn())
	credService := credential.NewService(credClient)

	binaryClient := client.NewBinaryClient(conn.Conn())
	binaryService := binary.NewService(binaryClient, fileStore, tokenProvider)

	return &Services{
		AuthSrv:       authService,
		CredsSrv:      credService,
		BinarySrv:     binaryService,
		TokenProvider: tokenProvider,
		grpcConn:      conn,
	}, nil
}

func (s *Services) Close() error {
	if s.grpcConn == nil {
		return fmt.Errorf("container: grpc connection is not initialized")
	}
	return s.grpcConn.Close()
}
