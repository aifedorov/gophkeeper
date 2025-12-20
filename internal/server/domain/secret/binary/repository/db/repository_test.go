package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockQuerier is a mock implementation of the Querier interface
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateFile(ctx context.Context, arg CreateFileParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) GetFile(ctx context.Context, arg GetFileParams) (File, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return File{}, args.Error(1)
	}
	return args.Get(0).(File), args.Error(1)
}

func (m *MockQuerier) ListFiles(ctx context.Context, userID uuid.UUID) ([]File, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]File), args.Error(1)
}

func (m *MockQuerier) DeleteFile(ctx context.Context, arg DeleteFileParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func TestRepository_Get(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		mockQuerier := new(MockQuerier)
		repo := &repository{
			queries: mockQuerier,
			logger:  logger,
		}

		userID := uuid.New()
		fileID := uuid.New()
		expectedFile := File{
			ID:             fileID,
			UserID:         userID,
			Name:           "test.txt",
			EncryptedPath:  []byte("encrypted-path"),
			EncryptedSize:  []byte("encrypted-size"),
			EncryptedNotes: []byte("encrypted-notes"),
			UploadedAt:     time.Now(),
		}

		mockQuerier.On("GetFile", ctx, GetFileParams{
			ID:     fileID,
			UserID: userID,
		}).Return(expectedFile, nil)

		result, err := repo.Get(ctx, userID.String(), fileID.String())

		require.NoError(t, err)
		assert.Equal(t, fileID.String(), result.ID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, expectedFile.Name, result.Name)
		assert.Equal(t, expectedFile.EncryptedPath, result.EncryptedPath)
		assert.Equal(t, expectedFile.EncryptedSize, result.EncryptedSize)
		assert.Equal(t, expectedFile.EncryptedNotes, result.EncryptedNotes)
		assert.Equal(t, expectedFile.UploadedAt, result.UploadedAt)
		mockQuerier.AssertExpectations(t)
	})

	t.Run("file not found", func(t *testing.T) {
		mockQuerier := new(MockQuerier)
		repo := &repository{
			queries: mockQuerier,
			logger:  logger,
		}

		userID := uuid.New()
		fileID := uuid.New()

		mockQuerier.On("GetFile", ctx, GetFileParams{
			ID:     fileID,
			UserID: userID,
		}).Return(nil, sql.ErrNoRows)

		result, err := repo.Get(ctx, userID.String(), fileID.String())

		require.Error(t, err)
		assert.ErrorIs(t, err, binary.ErrFileNotFound)
		assert.Equal(t, interfaces.RepositoryFile{}, result)
		mockQuerier.AssertExpectations(t)
	})

	t.Run("invalid user UUID", func(t *testing.T) {
		mockQuerier := new(MockQuerier)
		repo := &repository{
			queries: mockQuerier,
			logger:  logger,
		}

		fileID := uuid.New()

		result, err := repo.Get(ctx, "invalid-uuid", fileID.String())

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse user id")
		assert.Equal(t, interfaces.RepositoryFile{}, result)
		mockQuerier.AssertNotCalled(t, "GetFile")
	})

	t.Run("invalid file UUID", func(t *testing.T) {
		mockQuerier := new(MockQuerier)
		repo := &repository{
			queries: mockQuerier,
			logger:  logger,
		}

		userID := uuid.New()

		result, err := repo.Get(ctx, userID.String(), "invalid-uuid")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse file id")
		assert.Equal(t, interfaces.RepositoryFile{}, result)
		mockQuerier.AssertNotCalled(t, "GetFile")
	})

	t.Run("database error", func(t *testing.T) {
		mockQuerier := new(MockQuerier)
		repo := &repository{
			queries: mockQuerier,
			logger:  logger,
		}

		userID := uuid.New()
		fileID := uuid.New()
		dbErr := assert.AnError

		mockQuerier.On("GetFile", ctx, GetFileParams{
			ID:     fileID,
			UserID: userID,
		}).Return(nil, dbErr)

		result, err := repo.Get(ctx, userID.String(), fileID.String())

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get file")
		assert.Equal(t, interfaces.RepositoryFile{}, result)
		mockQuerier.AssertExpectations(t)
	})
}
