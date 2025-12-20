package repository

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type repository struct {
	queries Querier
	logger  *zap.Logger
}

func NewRepository(db DBTX, logger *zap.Logger) interfaces.Repository {
	return &repository{
		queries: New(db),
		logger:  logger,
	}
}

func (r *repository) Create(ctx context.Context, userID string, file interfaces.RepositoryFile) error {
	r.logger.Debug("repo: creating binary",
		zap.String("user_id", userID),
		zap.String("filename", file.Name),
	)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	fileUUID, err := uuid.Parse(file.ID)
	if err != nil {
		r.logger.Error("repo: failed to parse binary id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse binary id: %w", err)
	}

	r.logger.Debug("repo: IDs parsed successfully",
		zap.String("user_id", userID),
		zap.String("id", file.ID),
	)

	err = r.queries.CreateFile(ctx, CreateFileParams{
		ID:             fileUUID,
		UserID:         userUUID,
		Name:           file.Name,
		EncryptedPath:  file.EncryptedPath,
		EncryptedSize:  file.EncryptedSize,
		EncryptedNotes: file.EncryptedNotes,
		UploadedAt:     file.UploadedAt,
	})
	if err != nil {
		r.logger.Error("repo: failed to create binary", zap.Error(err))
		return fmt.Errorf("repo: failed to create binary: %w", err)
	}

	r.logger.Debug("repo: binary created successfully", zap.String("id", file.ID))
	return nil
}

func (r *repository) Get(ctx context.Context, userID, id string) (interfaces.RepositoryFile, error) {
	// TODO: implement
	panic("implement me")
}

func (r *repository) List(ctx context.Context, userID string) ([]interfaces.RepositoryFile, error) {
	r.logger.Debug("repo: listing files", zap.String("user_id", userID))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	dbFiles, err := r.queries.ListFiles(ctx, userUUID)
	if err != nil {
		r.logger.Error("repo: failed to list files", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to list files: %w", err)
	}

	files := make([]interfaces.RepositoryFile, len(dbFiles))
	for i, f := range dbFiles {
		files[i] = interfaces.RepositoryFile{
			ID:             f.ID.String(),
			UserID:         f.UserID.String(),
			Name:           f.Name,
			EncryptedPath:  f.EncryptedPath,
			EncryptedSize:  f.EncryptedSize,
			EncryptedNotes: f.EncryptedNotes,
			UploadedAt:     f.UploadedAt,
		}
	}
	return files, nil
}

func (r *repository) Delete(ctx context.Context, userID, id string) error {
	r.logger.Debug("repo: deleting binary", zap.String("user_id", userID), zap.String("id", id))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("repo: failed to parse binary id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse binary id: %w", err)
	}

	rows, err := r.queries.DeleteFile(ctx, DeleteFileParams{
		ID:     idUUID,
		UserID: userUUID,
	})
	if err != nil {
		r.logger.Error("repo: failed to delete binary", zap.Error(err))
		return fmt.Errorf("repo: failed to delete binary: %w", err)
	}
	if rows == 0 {
		r.logger.Debug("repo: binary not found")
		return binary.ErrFileNotFound
	}
	return nil
}
