// Package binary provides binary file management services for the GophKeeper server.
//
// This package implements the core business logic for managing binary file storage with
// end-to-end encryption. Files are encrypted using AES-256-GCM before storage and decrypted
// during retrieval. All operations require user authentication and use encryption keys
// derived from the user's password.
package binary

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/crypto"
	"go.uber.org/zap"
)

// Service defines the interface for binary file management operations.
// All methods require user authentication and encryption key for file encryption/decryption.
type Service interface {
	// Upload stores a new binary file for the specified user with encryption.
	// The file is read from the provided reader, encrypted using AES-256-GCM, and stored both
	// in the file storage and database. The encryptionKey should be provided as a base64-encoded string.
	// Returns ErrNameExists if a file with the same name already exists for the user.
	// If database creation fails, the uploaded file is automatically cleaned up.
	Upload(ctx context.Context, userID, encryptionKey string, metadata interfaces.FileMetadata, reader io.Reader) (*interfaces.File, error)
	// List retrieves all binary files for the specified user and decrypts their metadata.
	// The encryptionKey should be provided as a base64-encoded string.
	// Returns an empty slice if the user has no files.
	List(ctx context.Context, userID, encryptionKey string) ([]interfaces.File, error)
	// Download retrieves a binary file for the specified user and returns a reader for the decrypted content.
	// The encryptionKey should be provided as a base64-encoded string.
	// Returns ErrNotFound if the file doesn't exist or doesn't belong to the user.
	// The returned reader should be closed by the caller after use.
	Download(ctx context.Context, userID, encryptionKey, id string) (io.Reader, interfaces.FileMetadata, error)
	// Update updates a binary file for the specified user.
	// The file is updated in the database and file storage.
	// The encryptionKey should be provided as a base64-encoded string.
	// Returns an error if the file doesn't exist or doesn't belong to the user.
	Update(ctx context.Context, userID, encryptionKey string, metadata interfaces.FileMetadata, reader io.Reader) (*interfaces.File, error)
	// Delete removes a binary file for the specified user.
	// Deletes both the database record and the physical file.
	// Returns ErrNotFound if the file doesn't exist or doesn't belong to the user.
	Delete(ctx context.Context, userID, id string) error
}

// service implements the Service interface for binary file management.
type service struct {
	repo      interfaces.Repository
	fileStore interfaces.FileStorage
	crypto    interfaces.CryptoService
	logger    *zap.Logger
	locks     sync.Map // map[fileID]*sync.RWMutex
}

// NewService creates a new instance of the binary file service with the provided dependencies.
// It initializes the service with a file repository, file storage, cryptographic service, and logger.
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

// Upload stores a new binary file for the specified user with encryption.
// The file is read from the provided reader, encrypted using AES-256-GCM, and stored both
// in the file storage and database. The encryptionKey should be provided as a base64-encoded string.
// If database creation fails, the uploaded file is automatically cleaned up.
func (s *service) Upload(ctx context.Context, userID, encryptionKey string, metadata interfaces.FileMetadata, reader io.Reader) (*interfaces.File, error) {
	s.logger.Debug("binary: uploading file",
		zap.String("user_id", userID),
	)

	encryptReader, key, err := s.wrapToEncrypt(encryptionKey, reader)
	if err != nil {
		s.logger.Error("binary: failed to wrap reader for encryption", zap.Error(err))
		return nil, fmt.Errorf("failed to wrap reader for encryption: %w", err)
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

	removeFile := func() {
		err = s.fileStore.Delete(ctx, userID, file.GetID())
		if err != nil {
			s.logger.Error("binary: failed to delete uploaded file", zap.Error(err))
		}
	}

	s.logger.Debug("binary: file encrypted successfully", zap.String("filepath", filepath))

	repoFile, err := FileToRepository(s.crypto, key, file)
	if err != nil {
		s.logger.Error("binary: failed to convert file to repository", zap.Error(err))
		removeFile()
		return nil, fmt.Errorf("failed to convert file: %w", err)
	}

	f, err := s.repo.Create(ctx, userID, repoFile)
	if errors.Is(err, ErrNameExists) {
		s.logger.Debug("binary: file name already exists", zap.String("name", file.GetName()))
		removeFile()
		return nil, ErrNameExists
	}
	if err != nil {
		s.logger.Error("binary: failed to create file in repository", zap.Error(err))
		removeFile()
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	res, err := RepositoryToDomain(s.crypto, key, f)
	if err != nil {
		s.logger.Error("binary: failed to convert file metadata", zap.Error(err))
		removeFile()
		return nil, fmt.Errorf("failed to convert file metadata: %w", err)
	}

	s.logger.Debug("binary: file created successfully", zap.String("id", file.GetID()))
	return res, nil
}

// List retrieves all binary files for the specified user and decrypts their metadata.
// The encryptionKey should be provided as a base64-encoded string.
// Returns an empty slice if the user has no files.
func (s *service) List(ctx context.Context, userID, encryptionKey string) ([]interfaces.File, error) {
	s.logger.Debug("binary: listing files", zap.String("user_id", userID))

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil || len(key) != keyAES256Len {
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
		file, err := RepositoryToDomain(s.crypto, key, &repoFile)
		if err != nil || file == nil {
			s.logger.Error("binary: failed to convert file metadata", zap.Error(err))
			return nil, fmt.Errorf("failed to convert file metadata: %w", err)
		}
		res[i] = *file
	}
	return res, nil
}

// Download retrieves a binary file for the specified user and returns a reader for the decrypted content.
// The encryptionKey should be provided as a base64-encoded string.
// Returns ErrNotFound if the file doesn't exist or doesn't belong to the user.
// The returned reader should be closed by the caller after use.
// Uses a read lock to allow concurrent downloads of the same file.
func (s *service) Download(ctx context.Context, userID, encryptionKey, id string) (io.Reader, interfaces.FileMetadata, error) {
	s.logger.Debug("binary: downloading file", zap.String("user_id", userID), zap.String("id", id))

	mu := s.acquireReadLock(id)
	defer s.releaseReadLock(mu)

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil || len(key) != keyAES256Len {
		s.logger.Error("binary: failed to decode encryption key", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	fileRepo, err := s.repo.Get(ctx, userID, id)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("binary: file not found", zap.String("id", id))
		return nil, interfaces.FileMetadata{}, ErrNotFound
	}
	if err != nil {
		s.logger.Error("binary: failed to get file from repository", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to get file from repository: %w", err)
	}

	file, err := RepositoryToDomain(s.crypto, key, fileRepo)
	if err != nil {
		s.logger.Error("binary: failed to convert file metadata", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to convert file metadata: %w", err)
	}

	encryptReader, err := s.fileStore.Download(ctx, userID, file.GetID())
	if err != nil {
		s.logger.Error("binary: failed to open file for reading", zap.Error(err))
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to open file for reading: %w", err)
	}

	decryptingReader, err := crypto.NewDecryptReader(encryptReader, key)
	if err != nil {
		s.logger.Error("binary: failed to create decrypting reader", zap.Error(err))
		_ = encryptReader.Close()
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to create decrypting reader: %w", err)
	}

	meta, err := FileToMetadata(file)
	if err != nil {
		s.logger.Error("binary: failed to convert file metadata to domain", zap.Error(err))
		_ = encryptReader.Close()
		return nil, interfaces.FileMetadata{}, fmt.Errorf("failed to convert file metadata to domain: %w", err)
	}

	return decryptingReader, meta, nil
}

// Update updates a binary file for the specified user.
// The file is updated in the database and file storage with encryption.
// The encryptionKey should be provided as a base64-encoded string.
// Returns an error if the file doesn't exist or doesn't belong to the user.
// Returns ErrVersionConflict if the file version doesn't match (optimistic locking).
func (s *service) Update(ctx context.Context, userID, encryptionKey string, metadata interfaces.FileMetadata, reader io.Reader) (*interfaces.File, error) {
	s.logger.Debug("binary: updating file",
		zap.String("user_id", userID),
		zap.String("id", metadata.ID),
	)

	mu := s.acquireLock(metadata.ID)
	defer func() {
		s.releaseLock(mu)
	}()

	encryptReader, key, err := s.wrapToEncrypt(encryptionKey, reader)
	if err != nil {
		s.logger.Error("binary: failed to wrap reader for encryption", zap.Error(err))
		return nil, fmt.Errorf("failed to wrap reader for encryption: %w", err)
	}

	file, err := MetadataToFile(metadata)
	if err != nil {
		s.logger.Error("binary: failed to convert metadata to domain", zap.Error(err))
		return nil, fmt.Errorf("failed to convert metadata: %w", err)
	}

	s.logger.Debug("binary: begin updating file in file storage")

	tmppath, targetpath, err := s.fileStore.BeginUpdate(ctx, userID, file.GetID(), encryptReader)
	if err != nil {
		s.logger.Error("binary: failed to update file", zap.Error(err))
		_ = s.fileStore.AbortUpdate(ctx, tmppath)
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	file.SetPath(targetpath)

	s.logger.Debug("binary: file encrypted successfully",
		zap.String("tmppath", tmppath),
		zap.String("targetpath", targetpath),
	)

	s.logger.Debug("binary: updating file in repository")

	repoFile, err := FileToRepository(s.crypto, key, file)
	if err != nil {
		s.logger.Error("binary: failed to convert file to repository", zap.Error(err))
		_ = s.fileStore.AbortUpdate(ctx, tmppath)
		return nil, fmt.Errorf("failed to convert file: %w", err)
	}

	newRepoFile, err := s.repo.Update(ctx, userID, metadata.ID, repoFile)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("binary: file not found", zap.String("id", metadata.ID))
		_ = s.fileStore.AbortUpdate(ctx, tmppath)
		return nil, ErrNotFound
	}
	if errors.Is(err, ErrVersionConflict) {
		s.logger.Debug("binary: file version mismatch", zap.String("id", metadata.ID))
		_ = s.fileStore.AbortUpdate(ctx, tmppath)
		return nil, ErrVersionConflict
	}
	if err != nil {
		s.logger.Error("binary: failed to update file in repository", zap.Error(err))
		_ = s.fileStore.AbortUpdate(ctx, tmppath)
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	newFile, err := RepositoryToDomain(s.crypto, key, newRepoFile)
	if err != nil {
		s.logger.Error("binary: failed to convert repository file to domain", zap.Error(err))
		_ = s.fileStore.AbortUpdate(ctx, tmppath)
		return nil, fmt.Errorf("failed to convert repository file to domain: %w", err)
	}

	s.logger.Debug("binary: file updated in repository successfully", zap.String("id", metadata.ID))

	s.logger.Debug("binary: committing file update in file storage")
	err = s.fileStore.CommitUpdate(ctx, userID, newFile.GetID())
	if err != nil {
		s.logger.Error("binary: failed to commit file update", zap.Error(err))
		_ = s.fileStore.AbortUpdate(ctx, tmppath)
		return nil, fmt.Errorf("failed to commit file update: %w", err)
	}

	return newFile, nil
}

// Delete removes a binary file for the specified user.
// Deletes the database record first, then attempts to delete the physical file.
// If physical file deletion fails, it logs an error but doesn't return an error
// since the database record is already deleted.
// Returns ErrNotFound if the file doesn't exist or doesn't belong to the user.
// The lock for this file is removed from memory after successful deletion.
func (s *service) Delete(ctx context.Context, userID, id string) error {
	s.logger.Debug("binary: deleting file", zap.String("user_id", userID), zap.String("id", id))

	mu := s.acquireLock(id)
	defer func() {
		s.releaseLock(mu)
		s.removeLock(id)
	}()

	err := s.repo.Delete(ctx, userID, id)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("binary: file not found", zap.String("id", id))
		return ErrNotFound
	}
	if err != nil {
		s.logger.Error("binary: failed to delete file from repository", zap.Error(err))
		return fmt.Errorf("failed to delete file from repository: %w", err)
	}

	err = s.fileStore.Delete(ctx, userID, id)
	if err != nil {
		s.logger.Error("binary: failed to delete file from file storage", zap.Error(err))
	}

	s.logger.Debug("binary: file deleted successfully", zap.String("id", id))
	return nil
}
