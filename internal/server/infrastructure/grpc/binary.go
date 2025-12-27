package server

import (
	"context"
	"errors"
	"io"

	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/binary/v1"
	"github.com/aifedorov/gophkeeper/internal/server/config"
	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// bufferSize is the chunk size for streaming file data (1MB).
const bufferSize = 1024 * 1024

// BinaryServer implements the BinaryService gRPC service.
// It handles file upload, download, and management with encryption.
type BinaryServer struct {
	pb.UnimplementedBinaryServiceServer
	cfg       *config.Config
	logger    *zap.Logger
	authScr   auth.Service
	binarySrv binary.Service
}

// NewBinaryServer creates a new BinaryServer with the provided dependencies.
func NewBinaryServer(cfg *config.Config, logger *zap.Logger, authSrv auth.Service, binarySrv binary.Service) *BinaryServer {
	return &BinaryServer{
		cfg:       cfg,
		logger:    logger,
		authScr:   authSrv,
		binarySrv: binarySrv,
	}
}

// Upload receives a file stream and stores it encrypted.
// First message must contain metadata, subsequent messages contain chunks.
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
		return status.Error(codes.Internal, "internal error")
	}

	metadata := firstMsg.GetFile()
	if metadata == nil {
		s.logger.Error("grpc: file metadata is nil")
		return status.Error(codes.InvalidArgument, "invalid request")
	}

	streamReader := newUploadStreamReader(stream)
	fileMetadata := interfaces.FileMetadata{
		Name:  metadata.GetName(),
		Size:  metadata.GetSize(),
		Notes: metadata.GetNotes(),
	}

	res, err := s.binarySrv.Upload(ctx, userID, encryptionKey, fileMetadata, streamReader)
	if errors.Is(err, binary.ErrNameExists) {
		s.logger.Debug("grpc: file with this name already exists")
		return status.Error(codes.AlreadyExists, "file with this name already exists")
	}
	if err != nil {
		s.logger.Error("grpc: failed to upload file", zap.Error(err))
		return status.Errorf(codes.Internal, "internal binary error: %s", err.Error())
	}

	fileID := res.GetID()
	version := res.GetVersion()
	return stream.SendAndClose(&pb.UploadResponse{
		FileId:  &fileID,
		Version: &version,
	})
}

// List retrieves metadata for all files owned by the authenticated user.
func (s *BinaryServer) List(ctx context.Context, _ *pb.ListRequest) (*pb.ListResponse, error) {
	s.logger.Debug("grpc: list binary request received")

	userID, encryptionKey, err := s.authScr.GetUserDataFromContext(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	s.logger.Debug("grpc: user_id extracted from token", zap.String("user_id", userID))

	files, err := s.binarySrv.List(ctx, userID, encryptionKey)
	if err != nil {
		s.logger.Error("grpc: failed to get list of files", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal binary error")
	}

	s.logger.Debug("grpc: received list of files from repository", zap.Int("count", len(files)))

	resFiles := make([]*pb.ListResponse_Metadata, len(files))
	for i, f := range files {
		id := f.GetID()
		name := f.GetName()
		size := f.GetSize()
		notes := f.GetNotes()
		version := f.GetVersion()
		uploadedAt := f.GetUploadedAt()

		resFiles[i] = &pb.ListResponse_Metadata{
			Id:         &id,
			Name:       &name,
			Size:       &size,
			Notes:      &notes,
			Version:    &version,
			UploadedAt: timestamppb.New(uploadedAt),
		}
	}

	res := &pb.ListResponse{
		Files: resFiles,
	}
	return res, nil
}

// Download streams a file to the client.
// First response contains metadata, subsequent responses contain chunks.
func (s *BinaryServer) Download(req *pb.DownloadRequest, stream grpc.ServerStreamingServer[pb.DownloadResponse]) error {
	s.logger.Debug("grpc: download binary request received")
	ctx := stream.Context()

	s.logger.Debug("grpc: extracting user ID and encryption key from token")
	userID, encryptionKey, err := s.authScr.GetUserDataFromContext(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	s.logger.Debug("grpc: downloading binary")
	fieldID := req.GetFileId()
	if fieldID == "" {
		return status.Error(codes.InvalidArgument, "invalid request")
	}

	reader, meta, err := s.binarySrv.Download(ctx, userID, encryptionKey, fieldID)
	if err != nil {
		s.logger.Error("grpc: failed to download file", zap.Error(err))
		return status.Errorf(codes.Internal, "internal binary error: %s", err.Error())
	}

	err = stream.Send(&pb.DownloadResponse{
		Data: &pb.DownloadResponse_File{
			File: &pb.DownloadResponse_Metadata{
				Name:    &meta.Name,
				Size:    &meta.Size,
				Notes:   &meta.Notes,
				Version: &meta.Version,
			},
		},
	})
	if err != nil {
		return status.Errorf(codes.Internal, "internal binary error: %s", err.Error())
	}

	buf := make([]byte, bufferSize)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			err = stream.Send(&pb.DownloadResponse{
				Data: &pb.DownloadResponse_Chunk{
					Chunk: buf[:n],
				},
			})
			if err != nil {
				return status.Errorf(codes.Internal, "internal binary error: %s", err.Error())
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "internal binary error: %s", err.Error())
		}
	}
	return nil
}

// Delete removes a file by ID. Returns NotFound if file doesn't exist.
func (s *BinaryServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.logger.Debug("grpc: delete binary request received")

	userID, _, err := s.authScr.GetUserDataFromContext(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID from token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	fileID := req.GetFileId()
	if fileID == "" {
		return nil, status.Error(codes.InvalidArgument, "file_id is required")
	}

	err = s.binarySrv.Delete(ctx, userID, fileID)
	if err != nil {
		s.logger.Error("grpc: failed to delete file", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal binary error: %s", err.Error())
	}

	return &pb.DeleteResponse{}, nil
}

// Update replaces an existing file using optimistic locking.
// Returns NotFound if file doesn't exist, Aborted on version conflict.
func (s *BinaryServer) Update(stream grpc.ClientStreamingServer[pb.UpdateRequest, pb.UpdateResponse]) error {
	s.logger.Debug("grpc: update binary request received")

	ctx := stream.Context()

	s.logger.Debug("grpc: extracting user ID and encryption key from token")
	userID, encryptionKey, err := s.authScr.GetUserDataFromContext(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	s.logger.Debug("grpc: receiving metadata and uploading binary")

	firstMsg, err := stream.Recv()
	if err != nil {
		s.logger.Error("grpc: failed to receive first message", zap.Error(err))
		return status.Error(codes.Internal, "internal error")
	}

	metadata := firstMsg.GetFile()
	if metadata == nil {
		s.logger.Error("grpc: file metadata is nil")
		return status.Error(codes.InvalidArgument, "invalid request")
	}

	streamReader := newUpdateStreamReader(stream)
	fileMetadata := interfaces.FileMetadata{
		ID:      metadata.GetFileId(),
		Name:    metadata.GetName(),
		Size:    metadata.GetSize(),
		Notes:   metadata.GetNotes(),
		Version: metadata.GetVersion(),
	}

	res, err := s.binarySrv.Update(ctx, userID, encryptionKey, fileMetadata, streamReader)
	if errors.Is(err, binary.ErrNameExists) {
		s.logger.Debug("grpc: file with this name already exists")
		return status.Error(codes.AlreadyExists, "file with this name already exists")
	}
	if errors.Is(err, binary.ErrNotFound) {
		s.logger.Debug("grpc: file not found")
		return status.Error(codes.NotFound, "file not found")
	}
	if errors.Is(err, binary.ErrVersionConflict) {
		s.logger.Debug("grpc: file version conflict")
		return status.Error(codes.Aborted, "file was modified by another client, please refetch and retry")
	}
	if err != nil {
		s.logger.Error("grpc: failed to update file", zap.Error(err))
		return status.Errorf(codes.Internal, "internal binary error: %s", err.Error())
	}

	success := true
	version := res.GetVersion()
	return stream.SendAndClose(&pb.UpdateResponse{
		Success: &success,
		Version: &version,
	})
}
