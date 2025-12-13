package server

import (
	"context"
	"fmt"
	"net"

	authv1 "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/auth/v1"
	credv1 "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/credential/v1"
	"github.com/aifedorov/gophkeeper/internal/server/config"
	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/grpc/mw"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// GRPCServer defines the interface for running a gRPC server.
type GRPCServer interface {
	// Run starts the gRPC server and blocks until the server stops.
	// It sets up TLS credentials, registers services, and begins listening on the configured address.
	Run(ctx context.Context) error
}

type grpcServer struct {
	cfg        *config.Config
	logger     *zap.Logger
	grpc       *grpc.Server
	authServer *AuthServer
	credSrv    *CredentialServer
	jwtSrv     jwt.Service
	authSrv    auth.Service
}

// NewGRRPCServer creates a new instance of GRPCServer with the provided dependencies.
// It initializes the gRPC server that will handle authentication service requests.
func NewGRRPCServer(
	cfg *config.Config,
	logger *zap.Logger,
	grpc *grpc.Server,
	authServer *AuthServer,
	credSrv *CredentialServer,
	jwtSrv jwt.Service,
	authSrv auth.Service,
) GRPCServer {
	return &grpcServer{
		cfg:        cfg,
		logger:     logger,
		grpc:       grpc,
		authServer: authServer,
		credSrv:    credSrv,
		jwtSrv:     jwtSrv,
		authSrv:    authSrv,
	}
}

// Run starts the gRPC server and blocks until the server stops.
// It loads TLS credentials from certificate files, registers the AuthService,
// enables gRPC reflection, and begins listening on the configured address.
func (s *grpcServer) Run(ctx context.Context) error {
	s.logger.Info("starting grpc server", zap.String("addr", s.cfg.GRPCAddr))

	creds, err := credentials.NewServerTLSFromFile(
		"certs/server-cert.pem",
		"certs/server-key.pem",
	)
	if err != nil {
		return fmt.Errorf("failed to create credentials for grpc server: %w", err)
	}

	lc := net.ListenConfig{}
	listen, err := lc.Listen(ctx, "tcp", s.cfg.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	authInterceptor := mw.NewAuthInterceptor(s.jwtSrv, s.authSrv, s.logger)
	grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.UnaryAuthInterceptor),
	)

	s.grpc = grpc.NewServer(grpc.Creds(creds))
	authv1.RegisterAuthServiceServer(s.grpc, s.authServer)
	credv1.RegisterCredentialServiceServer(s.grpc, s.credSrv)

	reflection.Register(s.grpc)

	return s.grpc.Serve(listen)
}
