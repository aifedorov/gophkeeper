package binary

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/crypto"
	"go.uber.org/zap"
)

type Service interface {
	UploadFile(ctx context.Context, userID, encryptionKey string, metadata interfaces.FileMetadata, reader io.Reader) (*interfaces.File, error)
}

type service struct {
	repo      interfaces.Repository
	fileStore interfaces.FileStorage
	logger    *zap.Logger
}

func NewService(repo interfaces.Repository, fileStore interfaces.FileStorage, logger *zap.Logger) Service {
	return &service{
		repo:      repo,
		fileStore: fileStore,
		logger:    logger,
	}
}

func (s *service) UploadFile(ctx context.Context, userID, encryptionKey string, metadata interfaces.FileMetadata, reader io.Reader) (*interfaces.File, error) {
	s.logger.Debug("binary: uploading file",
		zap.String("user_id", userID),
		zap.String("filename", metadata.Name),
		zap.Int64("filesize", metadata.Size),
		zap.String("mimetype", metadata.MimeType),
	)

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		s.logger.Error("binary: failed to decode encryption key", zap.Error(err))
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	encryptReader, err := crypto.NewEncryptReader(reader, key)
	if err != nil {
		s.logger.Error("binary: failed to create encrypt reader", zap.Error(err))
		return nil, fmt.Errorf("failed to encrypt file: %w", err)
	}

	file, err := MetadataToDomain(metadata)
	if err != nil {
		s.logger.Error("binary: failed to convert metadata to domain", zap.Error(err))
		return nil, fmt.Errorf("failed to convert metadata: %w", err)
	}

	filepath, err := s.fileStore.Upload(ctx, userID, file.GetID(), encryptReader)
	if err != nil {
		s.logger.Error("binary: failed to upload file", zap.Error(err))
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	s.logger.Debug("binary: file encrypted successfully", zap.String("filepath", filepath))

	repoFile := DomainToRepository(file, filepath)
	err = s.repo.CreateFile(ctx, userID, repoFile)
	if err != nil {
		s.logger.Error("binary: failed to create file in repository", zap.Error(err))
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	s.logger.Debug("binary: file created successfully", zap.String("id", file.GetID()))
	return file, nil
}
