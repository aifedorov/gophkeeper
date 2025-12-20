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
	Upload(ctx context.Context, userID, encryptionKey string, metadata interfaces.FileMetadata, reader io.Reader) (*interfaces.File, error)
	List(ctx context.Context, userID, encryptionKey string) ([]interfaces.File, error)
}

type service struct {
	repo      interfaces.Repository
	fileStore interfaces.FileStorage
	crypto    interfaces.CryptoService
	logger    *zap.Logger
}

func NewService(
	repo interfaces.Repository,
	fileStore interfaces.FileStorage,
	crypto interfaces.CryptoService,
	logger *zap.Logger,
) Service {
	return &service{
		repo:      repo,
		fileStore: fileStore,
		crypto:    crypto,
		logger:    logger,
	}
}

func (s *service) Upload(ctx context.Context, userID, encryptionKey string, metadata interfaces.FileMetadata, reader io.Reader) (*interfaces.File, error) {
	s.logger.Debug("binary: uploading file",
		zap.String("user_id", userID),
		zap.String("filename", metadata.Name),
		zap.Int64("filesize", metadata.Size),
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

	s.logger.Debug("binary: converting file to repository file")
	repoFile, err := FileToRepository(s.crypto, key, file, filepath)
	if err != nil {
		s.logger.Error("binary: failed to convert file to repository", zap.Error(err))
		return nil, fmt.Errorf("failed to convert file: %w", err)
	}

	err = s.repo.Create(ctx, userID, repoFile)
	if err != nil {
		s.logger.Error("binary: failed to create file in repository", zap.Error(err))
		_ = s.fileStore.Delete(ctx, userID, file.GetID())
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	s.logger.Debug("binary: file created successfully", zap.String("id", file.GetID()))
	return file, nil
}

func (s *service) List(ctx context.Context, userID, encryptionKey string) ([]interfaces.File, error) {
	s.logger.Debug("binary: listing files", zap.String("user_id", userID))

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		s.logger.Error("binary: failed to decode encryption key", zap.Error(err))
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	repoFiles, err := s.repo.List(ctx, userID)
	if err != nil {
		s.logger.Error("binary: failed to list files", zap.Error(err))
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	res := make([]interfaces.File, len(repoFiles))
	for i, repoFile := range repoFiles {
		file, err := FileToDomain(s.crypto, key, repoFile)
		if err != nil || file == nil {
			s.logger.Error("binary: failed to convert file metadata", zap.Error(err))
			return nil, fmt.Errorf("failed to convert file metadata: %w", err)
		}
		res[i] = *file
	}
	return res, nil
}
