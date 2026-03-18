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

type repository struct {
	pool    TxBeginner
	queries Querier
	logger  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, db DBTX, logger *zap.Logger) interfaces.Repository {
	return &repository{
		pool:    pool,
		queries: newQuerier(db),
		logger:  logger,
	}
}

// newRepositoryForTest creates a repository with mocked dependencies for testing.
func newRepositoryForTest(pool TxBeginner, queries Querier, logger *zap.Logger) *repository {
	return &repository{
		pool:    pool,
		queries: queries,
		logger:  logger,
	}
}

func (r *repository) Create(ctx context.Context, userID string, file interfaces.RepositoryFile) (*interfaces.RepositoryFile, error) {
	r.logger.Debug("repo: creating binary",
		zap.String("user_id", userID),
		zap.String("filename", file.Name),
	)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	fileUUID, err := uuid.Parse(file.ID)
	if err != nil {
		r.logger.Error("repo: failed to parse binary id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse binary id: %w", err)
	}

	r.logger.Debug("repo: IDs parsed successfully",
		zap.String("user_id", userID),
		zap.String("id", file.ID),
	)

	f, err := r.queries.CreateFile(ctx, CreateFileParams{
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
		return nil, binary.ErrNameExists
	}
	if err != nil {
		r.logger.Error("repo: failed to create binary", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to create binary: %w", err)
	}

	r.logger.Debug("repo: binary created successfully", zap.String("id", file.ID))
	return &interfaces.RepositoryFile{
		ID:             f.ID.String(),
		UserID:         f.UserID.String(),
		Name:           f.Name,
		EncryptedPath:  f.EncryptedPath,
		EncryptedSize:  f.EncryptedSize,
		EncryptedNotes: f.EncryptedNotes,
		Version:        f.Version,
		UpdatedAt:      f.UpdatedAt,
	}, nil
}

func (r *repository) Get(ctx context.Context, userID, id string) (*interfaces.RepositoryFile, error) {
	r.logger.Debug("repo: getting file", zap.String("user_id", userID), zap.String("id", id))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("repo: failed to parse file id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse file id: %w", err)
	}

	dbFile, err := r.queries.GetFile(ctx, GetFileParams{
		ID:     idUUID,
		UserID: userUUID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		r.logger.Debug("repo: file not found")
		return nil, binary.ErrNotFound
	}
	if err != nil {
		r.logger.Error("repo: failed to get file", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to get file: %w", err)
	}

	r.logger.Debug("repo: file retrieved successfully", zap.String("id", id))
	return &interfaces.RepositoryFile{
		ID:             dbFile.ID.String(),
		UserID:         dbFile.UserID.String(),
		Name:           dbFile.Name,
		EncryptedPath:  dbFile.EncryptedPath,
		EncryptedSize:  dbFile.EncryptedSize,
		EncryptedNotes: dbFile.EncryptedNotes,
		Version:        dbFile.Version,
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
			Version:        f.Version,
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

func (r *repository) Update(ctx context.Context, userID, id string, in interfaces.RepositoryFile) (*interfaces.RepositoryFile, error) {
	r.logger.Debug("repo: updating binary", zap.String("user_id", userID), zap.String("id", id))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	r.logger.Debug("repo: start transaction for update")

	idUUID, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("repo: failed to parse binary id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse binary id: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.Error("repo: failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to begin transaction: %w", err)
	}
	defer func() {
		r.logger.Debug("repo: rollback transaction for update")
		_ = tx.Rollback(ctx)
	}()

	txQuery := r.queries.WithTx(tx)
	f, err := txQuery.GetFileForUpdate(ctx, GetFileForUpdateParams{
		ID:     idUUID,
		UserID: userUUID,
	})
	if notFoundError(err) {
		r.logger.Debug("repo: file not found for update", zap.String("id", id))
		return nil, binary.ErrNotFound
	}
	if err != nil {
		r.logger.Error("repo: failed to get in for update", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to get in for update: %w", err)
	}

	if f.Version != in.Version {
		r.logger.Debug("repo: version conflict",
			zap.Int64("clent_version", in.Version),
			zap.Int64("server_version", f.Version),
		)
		return nil, binary.ErrVersionConflict
	}

	f, err = txQuery.UpdateFile(ctx, UpdateFileParams{
		ID:             idUUID,
		UserID:         userUUID,
		Name:           in.Name,
		EncryptedPath:  in.EncryptedPath,
		EncryptedSize:  in.EncryptedSize,
		EncryptedNotes: in.EncryptedNotes,
	})
	if conflictError(err) {
		r.logger.Debug("repo: file name already exists", zap.String("name", in.Name))
		return nil, binary.ErrNameExists
	}
	if err != nil {
		r.logger.Error("repo: failed to update file", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to update file: %w", err)
	}

	r.logger.Debug("repo: in updated successfully", zap.String("id", id))
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: failed to commit transaction: %w", err)
	}

	return &interfaces.RepositoryFile{
		ID:             f.ID.String(),
		UserID:         f.UserID.String(),
		Name:           f.Name,
		EncryptedPath:  f.EncryptedPath,
		EncryptedSize:  f.EncryptedSize,
		EncryptedNotes: f.EncryptedNotes,
		Version:        f.Version,
		UpdatedAt:      f.UpdatedAt,
	}, nil
}
