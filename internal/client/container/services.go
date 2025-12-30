package container

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/config"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	authinterfaces "github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/aifedorov/gophkeeper/internal/client/domain/card"
	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	"github.com/aifedorov/gophkeeper/internal/client/domain/text"
	grpcClient "github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/grpc/client"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/storage/cache"
	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/storage/session"
	"github.com/aifedorov/gophkeeper/pkg/filestorage"
	"go.uber.org/zap"
)

type Services struct {
	AuthSrv       auth.Service
	CredsSrv      credential.Service
	BinarySrv     binary.Service
	CardSrv       card.Service
	TextSrv       text.Service
	TokenProvider authinterfaces.SessionProvider
	grpcConn      grpcClient.GRPCConnection
}

func NewServices(cfg *config.Config) (*Services, error) {
	sessionStore := session.NewStorage()
	cacheStore := cache.NewStorage()
	fileStore := filestorage.NewFileStorage(cfg.FileStoragePath, zap.NewNop())
	tokenProvider := auth.NewSessionProvider(sessionStore)
	conn, err := grpcClient.NewGRPCConnection(cfg.ServerAddr, tokenProvider)
	if err != nil {
		return nil, fmt.Errorf("container: failed to create grpc connection: %w", err)
	}

	authClient := client.NewAuthClient(conn.Conn())
	authService := auth.NewService(authClient, sessionStore)

	credClient := client.NewCredentialClient(conn.Conn())
	credService := credential.NewService(credClient, cacheStore)

	binaryClient := client.NewBinaryClient(conn.Conn())
	binaryService := binary.NewService(binaryClient, fileStore, cacheStore, tokenProvider)

	cardClient := client.NewCardClient(conn.Conn())
	cardService := card.NewService(cardClient, cacheStore)

	textService := text.NewService(binaryService, fileStore)

	return &Services{
		AuthSrv:       authService,
		CredsSrv:      credService,
		BinarySrv:     binaryService,
		CardSrv:       cardService,
		TextSrv:       textService,
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
