package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/auth/v1"
	"github.com/aifedorov/gophkeeper/internal/server/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer interface {
	Run(ctx context.Context) error
}

type grpcServer struct {
	cfg        *config.Config
	logger     *zap.Logger
	grpc       *grpc.Server
	authServer *AuthServer
}

func NewGRRPCServer(cfg *config.Config, logger *zap.Logger, grpc *grpc.Server, authServer *AuthServer) GRPCServer {
	return &grpcServer{
		cfg:        cfg,
		logger:     logger,
		grpc:       grpc,
		authServer: authServer,
	}
}

func (s *grpcServer) Run(ctx context.Context) error {
	s.logger.Info("starting grpc server", zap.String("addr", s.cfg.GRPCAddr))

	lc := net.ListenConfig{}
	listen, err := lc.Listen(ctx, "tcp", s.cfg.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.grpc = grpc.NewServer()
	pb.RegisterAuthServiceServer(s.grpc, s.authServer)

	reflection.Register(s.grpc)

	return s.grpc.Serve(listen)
}
