package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type testSetup struct {
	ctrl          *gomock.Controller
	mockPool      *MockTxBeginner
	mockQuerier   *MockQuerier
	mockTxQuerier *MockQuerier
	mockTx        *MockTx
	repo          *repository
	ctx           context.Context
	logger        *zap.Logger
	userID        uuid.UUID
	fileID        uuid.UUID
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	return &testSetup{
		ctrl:          ctrl,
		mockPool:      NewMockTxBeginner(ctrl),
		mockQuerier:   NewMockQuerier(ctrl),
		mockTxQuerier: NewMockQuerier(ctrl),
		mockTx:        NewMockTx(ctrl),
		ctx:           context.Background(),
		logger:        zap.NewNop(),
		userID:        uuid.New(),
		fileID:        uuid.New(),
	}
}

func (s *testSetup) initRepoForGet() {
	s.repo = newRepositoryForTest(nil, s.mockQuerier, s.logger)
}

func (s *testSetup) initRepoForUpdate() {
	s.repo = newRepositoryForTest(s.mockPool, s.mockQuerier, s.logger)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
}

func (s *testSetup) newTestFile() File {
	return File{
		ID:             s.fileID,
		UserID:         s.userID,
		Name:           "test_new.txt",
		EncryptedPath:  []byte("encrypted-path"),
		EncryptedSize:  []byte("encrypted-size"),
		EncryptedNotes: []byte("encrypted-notes"),
		UpdatedAt:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func (s *testSetup) newTestRepositoryFile() interfaces.RepositoryFile {
	return interfaces.RepositoryFile{
		Name:           "new-name.txt",
		EncryptedPath:  []byte("new-path"),
		EncryptedSize:  []byte("new-size"),
		EncryptedNotes: []byte("new-notes"),
	}
}

func (s *testSetup) expectGetSuccess(file File) {
	s.mockQuerier.EXPECT().
		GetFile(s.ctx, GetFileParams{ID: s.fileID, UserID: s.userID}).
		Return(file, nil).
		Times(1)
}

func (s *testSetup) expectGetNotFound() {
	s.mockQuerier.EXPECT().
		GetFile(s.ctx, GetFileParams{ID: s.fileID, UserID: s.userID}).
		Return(File{}, sql.ErrNoRows).
		Times(1)
}

func (s *testSetup) expectGetError(err error) {
	s.mockQuerier.EXPECT().
		GetFile(s.ctx, GetFileParams{ID: s.fileID, UserID: s.userID}).
		Return(File{}, err).
		Times(1)
}

func (s *testSetup) expectUpdateSuccess(existingFile File, newFile interfaces.RepositoryFile) {
	s.mockPool.EXPECT().Begin(s.ctx).Return(s.mockTx, nil).Times(1)
	s.mockQuerier.EXPECT().WithTx(s.mockTx).Return(s.mockTxQuerier).Times(1)
	s.mockTxQuerier.EXPECT().
		GetFileForUpdate(s.ctx, GetFileForUpdateParams{ID: s.fileID, UserID: s.userID}).
		Return(existingFile, nil).
		Times(1)
	s.mockTxQuerier.EXPECT().
		UpdateFile(s.ctx, UpdateFileParams{
			ID:             s.fileID,
			UserID:         s.userID,
			Name:           newFile.Name,
			EncryptedPath:  newFile.EncryptedPath,
			EncryptedSize:  newFile.EncryptedSize,
			EncryptedNotes: newFile.EncryptedNotes,
		}).
		Return(nil).
		Times(1)
	s.mockTx.EXPECT().Commit(s.ctx).Return(nil).Times(1)
	s.mockTx.EXPECT().Rollback(s.ctx).Return(nil).Times(1)
}

func (s *testSetup) expectUpdateNotFound() {
	s.mockPool.EXPECT().Begin(s.ctx).Return(s.mockTx, nil).Times(1)
	s.mockQuerier.EXPECT().WithTx(s.mockTx).Return(s.mockTxQuerier).Times(1)
	s.mockTxQuerier.EXPECT().
		GetFileForUpdate(s.ctx, GetFileForUpdateParams{ID: s.fileID, UserID: s.userID}).
		Return(File{}, sql.ErrNoRows).
		Times(1)
	s.mockTx.EXPECT().Rollback(s.ctx).Return(nil).Times(1)
}

func (s *testSetup) expectUpdateBeginError(err error) {
	s.mockPool.EXPECT().Begin(s.ctx).Return(nil, err).Times(1)
}

func (s *testSetup) expectUpdateFileError(existingFile File, err error) {
	s.mockPool.EXPECT().Begin(s.ctx).Return(s.mockTx, nil).Times(1)
	s.mockQuerier.EXPECT().WithTx(s.mockTx).Return(s.mockTxQuerier).Times(1)
	s.mockTxQuerier.EXPECT().
		GetFileForUpdate(s.ctx, GetFileForUpdateParams{ID: s.fileID, UserID: s.userID}).
		Return(existingFile, nil).
		Times(1)
	s.mockTxQuerier.EXPECT().
		UpdateFile(s.ctx, gomock.Any()).
		Return(err).
		Times(1)
	s.mockTx.EXPECT().Rollback(s.ctx).Return(nil).Times(1)
}

func (s *testSetup) expectUpdateCommitError(existingFile File) {
	s.mockPool.EXPECT().Begin(s.ctx).Return(s.mockTx, nil).Times(1)
	s.mockQuerier.EXPECT().WithTx(s.mockTx).Return(s.mockTxQuerier).Times(1)
	s.mockTxQuerier.EXPECT().
		GetFileForUpdate(s.ctx, GetFileForUpdateParams{ID: s.fileID, UserID: s.userID}).
		Return(existingFile, nil).
		Times(1)
	s.mockTxQuerier.EXPECT().
		UpdateFile(s.ctx, gomock.Any()).
		Return(nil).
		Times(1)
	s.mockTx.EXPECT().Commit(s.ctx).Return(errors.New("commit error")).Times(1)
	s.mockTx.EXPECT().Rollback(s.ctx).Return(nil).Times(1)
}

func assertError(t *testing.T, err error, wantErr bool, expectedErr error, errMsg string) {
	t.Helper()
	if wantErr {
		require.Error(t, err)
		if expectedErr != nil {
			assert.ErrorIs(t, err, expectedErr)
		}
		if errMsg != "" {
			assert.Contains(t, err.Error(), errMsg)
		}
	} else {
		require.NoError(t, err)
	}
}

func assertFileEqual(t *testing.T, got interfaces.RepositoryFile, want File) {
	t.Helper()
	assert.Equal(t, want.ID.String(), got.ID)
	assert.Equal(t, want.UserID.String(), got.UserID)
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.EncryptedPath, got.EncryptedPath)
	assert.Equal(t, want.EncryptedSize, got.EncryptedSize)
	assert.Equal(t, want.EncryptedNotes, got.EncryptedNotes)
	// Compare times with a small tolerance to account for execution time differences
	if !want.UpdatedAt.IsZero() && !got.UpdatedAt.IsZero() {
		assert.WithinDuration(t, want.UpdatedAt, got.UpdatedAt, time.Second)
	} else {
		assert.Equal(t, want.UpdatedAt, got.UpdatedAt)
	}
}
