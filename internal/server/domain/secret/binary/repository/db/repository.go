package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// QuerierFactory creates a Querier from a DBTX (used for transactions).
type QuerierFactory func(db DBTX) Querier

type repository struct {
	pool           TxBeginner
	queries        Querier
	querierFactory QuerierFactory
	logger         *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, db DBTX, logger *zap.Logger) interfaces.Repository {
	return &repository{
		pool:           pool,
		queries:        New(db),
		querierFactory: func(db DBTX) Querier { return New(db) },
		logger:         logger,
	}
}

// newRepositoryForTest creates a repository with mocked dependencies for testing.
func newRepositoryForTest(pool TxBeginner, queries Querier, txQueries Querier, logger *zap.Logger) *repository {
	return &repository{
		pool:           pool,
		queries:        queries,
		querierFactory: func(db DBTX) Querier { return txQueries },
		logger:         logger,
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
		UpdatedAt:      file.UpdatedAt,
	})
	if conflictError(err) {
		r.logger.Debug("repo: file name already exists", zap.String("name", file.Name))
		return binary.ErrNameExists
	}
	if err != nil {
		r.logger.Error("repo: failed to create binary", zap.Error(err))
		return fmt.Errorf("repo: failed to create binary: %w", err)
	}

	r.logger.Debug("repo: binary created successfully", zap.String("id", file.ID))
	return nil
}

func (r *repository) Get(ctx context.Context, userID, id string) (interfaces.RepositoryFile, error) {
	r.logger.Debug("repo: getting file", zap.String("user_id", userID), zap.String("id", id))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return interfaces.RepositoryFile{}, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("repo: failed to parse file id", zap.Error(err))
		return interfaces.RepositoryFile{}, fmt.Errorf("repo: failed to parse file id: %w", err)
	}

	dbFile, err := r.queries.GetFile(ctx, GetFileParams{
		ID:     idUUID,
		UserID: userUUID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		r.logger.Debug("repo: file not found")
		return interfaces.RepositoryFile{}, binary.ErrNotFound
	}
	if err != nil {
		r.logger.Error("repo: failed to get file", zap.Error(err))
		return interfaces.RepositoryFile{}, fmt.Errorf("repo: failed to get file: %w", err)
	}

	r.logger.Debug("repo: file retrieved successfully", zap.String("id", id))
	return interfaces.RepositoryFile{
		ID:             dbFile.ID.String(),
		UserID:         dbFile.UserID.String(),
		Name:           dbFile.Name,
		EncryptedPath:  dbFile.EncryptedPath,
		EncryptedSize:  dbFile.EncryptedSize,
		EncryptedNotes: dbFile.EncryptedNotes,
		UpdatedAt:      dbFile.UpdatedAt,
	}, nil
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
			UpdatedAt:      f.UpdatedAt,
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
		return binary.ErrNotFound
	}
	return nil
}

func (r *repository) Update(ctx context.Context, userID, id string, file interfaces.RepositoryFile) error {
	r.logger.Debug("repo: updating binary", zap.String("user_id", userID), zap.String("id", id))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	r.logger.Debug("repo: start transaction for update")

	idUUID, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("repo: failed to parse binary id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse binary id: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.Error("repo: failed to begin transaction", zap.Error(err))
		return fmt.Errorf("repo: failed to begin transaction: %w", err)
	}
	defer func() {
		r.logger.Debug("repo: rollback transaction for update")
		_ = tx.Rollback(ctx)
	}()

	txQuery := r.querierFactory(tx)
	_, err = txQuery.GetFileForUpdate(ctx, GetFileForUpdateParams{
		ID:     idUUID,
		UserID: userUUID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		r.logger.Debug("repo: file not found for update")
		return binary.ErrNotFound
	}
	if err != nil {
		r.logger.Error("repo: failed to get file for update", zap.Error(err))
		return fmt.Errorf("repo: failed to get file for update: %w", err)
	}

	err = txQuery.UpdateFile(ctx, UpdateFileParams{
		ID:             idUUID,
		UserID:         userUUID,
		Name:           file.Name,
		EncryptedPath:  file.EncryptedPath,
		EncryptedSize:  file.EncryptedSize,
		EncryptedNotes: file.EncryptedNotes,
	})
	if err != nil {
		r.logger.Error("repo: failed to update file", zap.Error(err))
		return fmt.Errorf("repo: failed to update file: %w", err)
	}
	r.logger.Debug("repo: file updated successfully", zap.String("id", id))
	return tx.Commit(ctx)
}
