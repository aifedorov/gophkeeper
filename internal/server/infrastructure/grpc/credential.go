package server

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/credential/v1"
	"github.com/aifedorov/gophkeeper/internal/server/config"
	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CredentialServer implements the CredentialService gRPC service.
// It handles CRUD operations for user credentials with encryption support.
// All operations require authentication via JWT token in request metadata.
type CredentialServer struct {
	pb.UnimplementedCredentialServiceServer
	cfg     *config.Config
	logger  *zap.Logger
	authSev auth.Service
	credSrv credential.Service
}

// NewCredentialServer creates a new instance of CredentialServer with the provided dependencies.
// It initializes the gRPC credential server that handles credential CRUD operations.
func NewCredentialServer(cfg *config.Config, logger *zap.Logger, authSev auth.Service, credSrv credential.Service) *CredentialServer {
	return &CredentialServer{
		cfg:     cfg,
		logger:  logger,
		authSev: authSev,
		credSrv: credSrv,
	}
}

// Create handles credential creation requests.
// It validates the request, extracts user ID from JWT token, creates the credential entity,
// and stores it with encryption. Returns the ID of the newly created credential.
// Returns Unauthenticated if JWT token is invalid, InvalidArgument if request data is invalid,
// or Internal if an unexpected error occurs.
func (s *CredentialServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	s.logger.Debug("grpc: create credential request received", zap.String("name", req.GetName()))

	userID, encryptionKey, err := s.authUser(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	newCred, err := credential.NewCredential(req.GetName(), req.GetLogin(), req.GetPassword(), req.GetNotes())
	if err != nil || newCred == nil {
		s.logger.Error("grpc: failed to create credential entity", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	res, err := s.credSrv.Create(ctx, userID, encryptionKey, *newCred)
	if errors.Is(err, credential.ErrNameExists) {
		s.logger.Debug("grpc: credential name already exists", zap.String("name", newCred.GetName()))
		return nil, status.Error(codes.AlreadyExists, "credential name already exists")
	}
	if err != nil {
		s.logger.Error("grpc: failed to create credential", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	id := res.GetID()
	s.logger.Debug("grpc: credential created successfully", zap.String("id", id))

	resp := pb.CreateResponse{
		Id: &id,
	}

	return &resp, nil
}

// List handles credential list requests.
// It validates the request, extracts user ID and encryption key from JWT token,
// retrieves the list of credentials, and returns the list of credentials.
// Returns Unauthenticated if JWT token is invalid, or Internal if an unexpected error occurs.
func (s *CredentialServer) List(ctx context.Context, _ *pb.ListRequest) (*pb.ListResponse, error) {
	s.logger.Debug("grpc: list credential request received")

	userID, encryptionKey, err := s.authUser(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	creds, err := s.credSrv.List(ctx, userID, encryptionKey)
	if err != nil {
		s.logger.Error("grpc: failed to get list of credentials", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	credentials := make([]*pb.ListResponse_ListItem, len(creds))
	for i, cred := range creds {
		id := cred.GetID()
		name := cred.GetName()
		login := cred.GetLogin()
		password := cred.GetPassword()
		notes := cred.GetMetadata()

		credentials[i] = &pb.ListResponse_ListItem{
			Id:       &id,
			Name:     &name,
			Login:    &login,
			Password: &password,
			Notes:    &notes,
		}
	}

	return &pb.ListResponse{
		Credentials: credentials,
	}, nil
}

func (s *CredentialServer) authUser(ctx context.Context) (userID string, encryptionKey string, err error) {
	userID, err = s.authSev.GetUserIDFromContext(ctx)
	if err != nil {
		return "", "", fmt.Errorf("grpc: failed to get userID: %w", err)
	}

	encryptionKey, err = s.authSev.GetEncryptionKeyFromContext(ctx)
	if err != nil {
		return "", "", fmt.Errorf("grpc: failed to get encryption key: %w", err)
	}

	return userID, encryptionKey, nil
}
