package server

import (
	"context"

	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/credential/v1"
	"github.com/aifedorov/gophkeeper/internal/server/config"
	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CredentialService struct {
	pb.UnimplementedCredentialServiceServer
	cfg     *config.Config
	logger  *zap.Logger
	authSev auth.Service
	credSrv credential.Service
}

func NewService(cfg *config.Config, logger *zap.Logger, authSev auth.Service, credSrv credential.Service) *CredentialService {
	return &CredentialService{
		cfg:     cfg,
		logger:  logger,
		authSev: authSev,
		credSrv: credSrv,
	}
}

func (s *CredentialService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	s.logger.Debug("Create called", zap.Any("req", req))

	userIDString, err := s.authSev.GetUserIDFromContext(ctx)
	if err != nil {
		s.logger.Error("Failed to get userID", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		s.logger.Error("Failed to parse user id", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid request")
	}

	newCred, err := credential.NewCredential(req.GetName(), req.GetLogin(), req.GetPassword(), req.GetMetadata())
	if err != nil || newCred == nil {
		s.logger.Error("Failed to create credential", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	res, err := s.credSrv.Create(ctx, userID, *newCred)
	if err != nil {
		s.logger.Error("Failed to create credential", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	id := res.GetID().String()
	resp := pb.CreateResponse{
		Id: &id,
	}

	return &resp, nil
}
