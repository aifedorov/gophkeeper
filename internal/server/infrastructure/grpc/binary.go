package server

import (
	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/binary/v1"
	"github.com/aifedorov/gophkeeper/internal/server/config"
	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BinaryServer struct {
	pb.UnimplementedBinaryServiceServer
	cfg       *config.Config
	logger    *zap.Logger
	authScr   auth.Service
	binarySrv binary.Service
}

func NewBinaryServer(cfg *config.Config, logger *zap.Logger, authSrv auth.Service, binarySrv binary.Service) *BinaryServer {
	return &BinaryServer{
		cfg:       cfg,
		logger:    logger,
		authScr:   authSrv,
		binarySrv: binarySrv,
	}
}

func (s *BinaryServer) Upload(stream grpc.ClientStreamingServer[pb.UploadRequest, pb.UploadResponse]) error {
	s.logger.Debug("grpc: upload binary request received")
	ctx := stream.Context()

	s.logger.Debug("grpc: extracting user ID and encryption key from token")
	userID, encryptionKey, err := s.authScr.GetUserDataFromContext(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	s.logger.Debug("grpc: uploading binary")
	firstMsg, err := stream.Recv()
	if err != nil {
		s.logger.Error("grpc: failed to receive first message", zap.Error(err))
		return status.Error(codes.Internal, "internal binary error")
	}

	metadata := firstMsg.GetFile()
	if metadata == nil {
		s.logger.Error("grpc: file metadata is nil")
		return status.Error(codes.InvalidArgument, "invalid request")
	}

	streamReader := newGRPCStreamReader(stream)
	fileMetadata := interfaces.FileMetadata{
		Name:     metadata.GetName(),
		Size:     metadata.GetSize(),
		MimeType: metadata.GetMimeType(),
	}

	res, err := s.binarySrv.UploadFile(ctx, userID, encryptionKey, fileMetadata, streamReader)
	if err != nil {
		s.logger.Error("grpc: failed to upload file", zap.Error(err))
		return status.Errorf(codes.Internal, "internal binary error: %s", err.Error())
	}

	fileID := res.GetID()
	return stream.SendAndClose(&pb.UploadResponse{
		FileId: &fileID,
	})
}
