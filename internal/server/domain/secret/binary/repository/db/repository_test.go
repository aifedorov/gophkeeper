package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestRepository_Get(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := newRepositoryForTest(nil, mockQuerier, nil, logger)

		userID := uuid.New()
		fileID := uuid.New()
		expectedFile := File{
			ID:             fileID,
			UserID:         userID,
			Name:           "test_new.txt",
			EncryptedPath:  []byte("encrypted-path"),
			EncryptedSize:  []byte("encrypted-size"),
			EncryptedNotes: []byte("encrypted-notes"),
			UpdatedAt:      time.Now(),
		}

		mockQuerier.EXPECT().
			GetFile(ctx, GetFileParams{ID: fileID, UserID: userID}).
			Return(expectedFile, nil)

		result, err := repo.Get(ctx, userID.String(), fileID.String())

		require.NoError(t, err)
		assert.Equal(t, fileID.String(), result.ID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, expectedFile.Name, result.Name)
		assert.Equal(t, expectedFile.EncryptedPath, result.EncryptedPath)
		assert.Equal(t, expectedFile.EncryptedSize, result.EncryptedSize)
		assert.Equal(t, expectedFile.EncryptedNotes, result.EncryptedNotes)
		assert.Equal(t, expectedFile.UpdatedAt, result.UpdatedAt)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := newRepositoryForTest(nil, mockQuerier, nil, logger)

		userID := uuid.New()
		fileID := uuid.New()

		mockQuerier.EXPECT().
			GetFile(ctx, GetFileParams{ID: fileID, UserID: userID}).
			Return(File{}, sql.ErrNoRows)

		result, err := repo.Get(ctx, userID.String(), fileID.String())

		require.Error(t, err)
		assert.ErrorIs(t, err, binary.ErrNotFound)
		assert.Equal(t, interfaces.RepositoryFile{}, result)
	})

	t.Run("invalid user UUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := newRepositoryForTest(nil, mockQuerier, nil, logger)

		fileID := uuid.New()

		result, err := repo.Get(ctx, "invalid-uuid", fileID.String())

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse user id")
		assert.Equal(t, interfaces.RepositoryFile{}, result)
	})

	t.Run("invalid file UUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := newRepositoryForTest(nil, mockQuerier, nil, logger)

		userID := uuid.New()

		result, err := repo.Get(ctx, userID.String(), "invalid-uuid")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse file id")
		assert.Equal(t, interfaces.RepositoryFile{}, result)
	})

	t.Run("database error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := newRepositoryForTest(nil, mockQuerier, nil, logger)

		userID := uuid.New()
		fileID := uuid.New()
		dbErr := errors.New("database error")

		mockQuerier.EXPECT().
			GetFile(ctx, GetFileParams{ID: fileID, UserID: userID}).
			Return(File{}, dbErr)

		result, err := repo.Get(ctx, userID.String(), fileID.String())

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get file")
		assert.Equal(t, interfaces.RepositoryFile{}, result)
	})
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := NewMockTxBeginner(ctrl)
		mockQuerier := NewMockQuerier(ctrl)
		mockTxQuerier := NewMockQuerier(ctrl)
		tx := NewMockTx(ctrl)

		repo := newRepositoryForTest(mockPool, mockQuerier, mockTxQuerier, logger)

		userID := uuid.New()
		fileID := uuid.New()
		existingFile := File{
			ID:             fileID,
			UserID:         userID,
			Name:           "test_new.txt",
			EncryptedPath:  []byte("encrypted-path"),
			EncryptedSize:  []byte("encrypted-size"),
			EncryptedNotes: []byte("encrypted-notes"),
			UpdatedAt:      time.Now(),
		}

		newFile := interfaces.RepositoryFile{
			Name:           "new-name.txt",
			EncryptedPath:  []byte("new-path"),
			EncryptedSize:  []byte("new-size"),
			EncryptedNotes: []byte("new-notes"),
		}

		mockPool.EXPECT().Begin(ctx).Return(tx, nil)
		mockTxQuerier.EXPECT().
			GetFileForUpdate(ctx, GetFileForUpdateParams{ID: fileID, UserID: userID}).
			Return(existingFile, nil)
		mockTxQuerier.EXPECT().
			UpdateFile(ctx, UpdateFileParams{
				ID:             fileID,
				UserID:         userID,
				Name:           newFile.Name,
				EncryptedPath:  newFile.EncryptedPath,
				EncryptedSize:  newFile.EncryptedSize,
				EncryptedNotes: newFile.EncryptedNotes,
			}).
			Return(nil)
		tx.EXPECT().Commit(ctx).Return(nil)
		tx.EXPECT().Rollback(ctx).Return(nil)

		err := repo.Update(ctx, userID.String(), fileID.String(), newFile)

		require.NoError(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := NewMockTxBeginner(ctrl)
		mockQuerier := NewMockQuerier(ctrl)
		mockTxQuerier := NewMockQuerier(ctrl)
		tx := NewMockTx(ctrl)

		repo := newRepositoryForTest(mockPool, mockQuerier, mockTxQuerier, logger)

		userID := uuid.New()
		fileID := uuid.New()

		mockPool.EXPECT().Begin(ctx).Return(tx, nil)
		mockTxQuerier.EXPECT().
			GetFileForUpdate(ctx, GetFileForUpdateParams{ID: fileID, UserID: userID}).
			Return(File{}, sql.ErrNoRows)
		tx.EXPECT().Rollback(ctx).Return(nil)

		err := repo.Update(ctx, userID.String(), fileID.String(), interfaces.RepositoryFile{})

		require.Error(t, err)
		assert.ErrorIs(t, err, binary.ErrNotFound)
	})

	t.Run("begin transaction error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := NewMockTxBeginner(ctrl)
		mockQuerier := NewMockQuerier(ctrl)

		repo := newRepositoryForTest(mockPool, mockQuerier, nil, logger)

		userID := uuid.New()
		fileID := uuid.New()
		txErr := errors.New("transaction error")

		mockPool.EXPECT().Begin(ctx).Return(nil, txErr)

		err := repo.Update(ctx, userID.String(), fileID.String(), interfaces.RepositoryFile{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to begin transaction")
	})

	t.Run("update file error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := NewMockTxBeginner(ctrl)
		mockQuerier := NewMockQuerier(ctrl)
		mockTxQuerier := NewMockQuerier(ctrl)
		tx := NewMockTx(ctrl)

		repo := newRepositoryForTest(mockPool, mockQuerier, mockTxQuerier, logger)

		userID := uuid.New()
		fileID := uuid.New()
		existingFile := File{
			ID:     fileID,
			UserID: userID,
			Name:   "test_new.txt",
		}
		updateErr := errors.New("update error")

		mockPool.EXPECT().Begin(ctx).Return(tx, nil)
		mockTxQuerier.EXPECT().
			GetFileForUpdate(ctx, GetFileForUpdateParams{ID: fileID, UserID: userID}).
			Return(existingFile, nil)
		mockTxQuerier.EXPECT().
			UpdateFile(ctx, gomock.Any()).
			Return(updateErr)
		tx.EXPECT().Rollback(ctx).Return(nil)

		err := repo.Update(ctx, userID.String(), fileID.String(), interfaces.RepositoryFile{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update file")
	})

	t.Run("invalid user UUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := newRepositoryForTest(nil, mockQuerier, nil, logger)

		fileID := uuid.New()

		err := repo.Update(ctx, "invalid-uuid", fileID.String(), interfaces.RepositoryFile{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse user id")
	})

	t.Run("invalid file UUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := newRepositoryForTest(nil, mockQuerier, nil, logger)

		userID := uuid.New()

		err := repo.Update(ctx, userID.String(), "invalid-uuid", interfaces.RepositoryFile{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse binary id")
	})

	t.Run("commit error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := NewMockTxBeginner(ctrl)
		mockQuerier := NewMockQuerier(ctrl)
		mockTxQuerier := NewMockQuerier(ctrl)
		tx := NewMockTx(ctrl)

		repo := newRepositoryForTest(mockPool, mockQuerier, mockTxQuerier, logger)

		userID := uuid.New()
		fileID := uuid.New()
		existingFile := File{
			ID:     fileID,
			UserID: userID,
			Name:   "test_new.txt",
		}

		mockPool.EXPECT().Begin(ctx).Return(tx, nil)
		mockTxQuerier.EXPECT().
			GetFileForUpdate(ctx, GetFileForUpdateParams{ID: fileID, UserID: userID}).
			Return(existingFile, nil)
		mockTxQuerier.EXPECT().
			UpdateFile(ctx, gomock.Any()).
			Return(nil)
		tx.EXPECT().Commit(ctx).Return(errors.New("commit error"))
		tx.EXPECT().Rollback(ctx).Return(nil)

		err := repo.Update(ctx, userID.String(), fileID.String(), interfaces.RepositoryFile{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "commit error")
	})
}
