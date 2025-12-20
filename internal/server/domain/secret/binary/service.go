package binary

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/crypto"
	"go.uber.org/zap"
)

type Service interface {
	Upload(ctx context.Context, userID, encryptionKey string, metadata interfaces.FileMetadata, reader io.Reader) (*interfaces.File, error)
	List(ctx context.Context, userID, encryptionKey string) ([]interfaces.File, error)
	Download(ctx context.Context, userID, encryptionKey, id string) (io.Reader, interfaces.FileMetadata, error)
	Delete(ctx context.Context, userID, id string) error
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

	file, err := MetadataToFile(metadata)
	if err != nil {
		s.logger.Error("binary: failed to convert metadata to domain", zap.Error(err))
		return nil, fmt.Errorf("failed to convert metadata: %w", err)
	}

	filepath, err := s.fileStore.Upload(ctx, userID, file.GetID(), encryptReader)
	if err != nil {
		s.logger.Error("binary: failed to upload file", zap.Error(err))
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	file.SetPath(filepath)

	s.logger.Debug("binary: file encrypted successfully", zap.String("filepath", filepath))

	s.logger.Debug("binary: converting file to repository file")
	repoFile, err := FileToRepository(s.crypto, key, file)
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

func (s *service) Download(ctx context.Context, userID, encryptionKey, id string) (io.Reader, interfaces.FileMetadata, error) {
	s.logger.Debug("binary: downloading file", zap.String("user_id", userID), zap.String("id", id))

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		s.logger.Error("binary: failed to decode encryption key", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	fileRepo, err := s.repo.Get(ctx, userID, id)
	if errors.Is(err, ErrFileNotFound) {
		s.logger.Debug("binary: file not found", zap.String("id", id))
		return nil, interfaces.FileMetadata{}, ErrFileNotFound
	}
	if err != nil {
		s.logger.Error("binary: failed to get file from repository", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to get file from repository: %w", err)
	}

	file, err := FileToDomain(s.crypto, key, fileRepo)
	if err != nil {
		s.logger.Error("binary: failed to convert file metadata", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to convert file metadata: %w", err)
	}

	encryptedReader, err := s.fileStore.Download(ctx, userID, file.GetID())
	if err != nil {
		s.logger.Error("binary: failed to open file for reading", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to open file for reading: %w", err)
	}

	decryptingReader, err := crypto.NewDecryptReader(encryptedReader, key)
	if err != nil {
		s.logger.Error("binary: failed to create decrypting reader", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to create decrypting reader: %w", err)
	}

	meta, err := FileToMetadata(file)
	if err != nil {
		s.logger.Error("binary: failed to convert file metadata to domain", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to convert file metadata to domain: %w", err)
	}

	return decryptingReader, meta, nil
}

func (s *service) Delete(ctx context.Context, userID, id string) error {
	s.logger.Debug("binary: deleting file", zap.String("user_id", userID), zap.String("id", id))

	// Delete from repository first
	err := s.repo.Delete(ctx, userID, id)
	if errors.Is(err, ErrFileNotFound) {
		s.logger.Debug("binary: file not found", zap.String("id", id))
		return ErrFileNotFound
	}
	if err != nil {
		s.logger.Error("binary: failed to delete file from repository", zap.Error(err))
		return fmt.Errorf("failed to delete file from repository: %w", err)
	}

	// Delete physical file
	err = s.fileStore.Delete(ctx, userID, id)
	if err != nil {
		s.logger.Warn("binary: failed to delete physical file", zap.Error(err))
		// Don't return error here as the database record is already deleted
	}

	s.logger.Debug("binary: file deleted successfully", zap.String("id", id))
	return nil
}
