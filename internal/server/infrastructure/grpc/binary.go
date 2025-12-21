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

const bufferSize = 1024 * 1024

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
	return stream.SendAndClose(&pb.UploadResponse{
		FileId: &fileID,
	})
}

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

	resFiles := make([]*pb.MetadataResponse, len(files))
	for i, f := range files {
		id := f.GetID()
		name := f.GetName()
		size := f.GetSize()
		notes := f.GetNotes()
		uploadedAt := f.GetUploadedAt()

		resFiles[i] = &pb.MetadataResponse{
			Id:         &id,
			Name:       &name,
			Size:       &size,
			Notes:      &notes,
			UploadedAt: timestamppb.New(uploadedAt),
		}
	}

	res := &pb.ListResponse{
		Files: resFiles,
	}
	return res, nil
}

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
				Name: &meta.Name,
				Size: &meta.Size,
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
		ID:    metadata.GetFileId(),
		Name:  metadata.GetName(),
		Size:  metadata.GetSize(),
		Notes: metadata.GetNotes(),
	}

	_, err = s.binarySrv.Update(ctx, userID, encryptionKey, fileMetadata, streamReader)
	if err != nil {
		s.logger.Error("grpc: failed to upload file", zap.Error(err))
		return status.Errorf(codes.Internal, "internal binary error: %s", err.Error())
	}

	success := true
	return stream.SendAndClose(&pb.UpdateResponse{
		Success: &success,
	})
}
