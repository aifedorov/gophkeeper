package binary

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	authinterfaces "github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/pkg/filestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockClient is a mock implementation of Client interface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Upload(ctx context.Context, fileInfo *FileInfo, reader io.Reader) error {
	args := m.Called(ctx, fileInfo, reader)
	return args.Error(0)
}

func (m *MockClient) List(ctx context.Context) ([]File, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]File), args.Error(1)
}

func (m *MockClient) Download(ctx context.Context, id string) (io.ReadCloser, *FileMeta, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(io.ReadCloser), nil, args.Error(2)
	}
	return args.Get(0).(io.ReadCloser), args.Get(1).(*FileMeta), args.Error(2)
}

func (m *MockClient) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSessionProvider is a mock implementation of SessionProvider
type MockSessionProvider struct {
	mock.Mock
}

func (m *MockSessionProvider) GetSession(ctx context.Context) (authinterfaces.Session, error) {
	args := m.Called(ctx)
	return args.Get(0).(authinterfaces.Session), args.Error(1)
}

// mockReadCloser is a simple implementation of io.ReadCloser for testing
type mockReadCloser struct {
	io.Reader
	closed bool
}

func (m *mockReadCloser) Close() error {
	m.closed = true
	return nil
}

func TestService_Download(t *testing.T) {
	t.Parallel()

	t.Run("successful download", func(t *testing.T) {
		mockClient := new(MockClient)
		mockSessionProvider := new(MockSessionProvider)
		storage := filestorage.NewFileStorage(zap.NewNop())

		service := &service{
			client:          mockClient,
			store:           storage,
			sessionProvider: mockSessionProvider,
		}

		ctx := context.Background()
		fileID := "test-file-id"
		userID := "test-user-id"

		session := authinterfaces.NewSession("token", "key", userID)
		mockSessionProvider.On("GetSession", ctx).Return(session, nil)

		fileContent := "test file content"
		mockReader := &mockReadCloser{Reader: strings.NewReader(fileContent)}
		mockMeta := &FileMeta{
			name:  "test.txt",
			size:  int64(len(fileContent)),
			notes: "test notes",
		}

		mockClient.On("Download", ctx, fileID).Return(mockReader, mockMeta, nil)

		filepath, err := service.Download(ctx, fileID)

		require.NoError(t, err)
		assert.NotEmpty(t, filepath)
		assert.Contains(t, filepath, "test.txt", "filepath should contain the real filename")
		assert.True(t, mockReader.closed, "reader should be closed")
		mockClient.AssertExpectations(t)
		mockSessionProvider.AssertExpectations(t)
	})

	t.Run("session provider error", func(t *testing.T) {
		mockClient := new(MockClient)
		mockSessionProvider := new(MockSessionProvider)
		storage := filestorage.NewFileStorage(zap.NewNop())

		service := &service{
			client:          mockClient,
			store:           storage,
			sessionProvider: mockSessionProvider,
		}

		ctx := context.Background()
		fileID := "test-file-id"

		sessionErr := errors.New("session not found")
		mockSessionProvider.On("GetSession", ctx).Return(authinterfaces.Session{}, sessionErr)

		filepath, err := service.Download(ctx, fileID)

		require.Error(t, err)
		assert.Empty(t, filepath)
		assert.Contains(t, err.Error(), "failed to get session")
		mockSessionProvider.AssertExpectations(t)
		mockClient.AssertNotCalled(t, "Download")
	})

	t.Run("client download error", func(t *testing.T) {
		mockClient := new(MockClient)
		mockSessionProvider := new(MockSessionProvider)
		storage := filestorage.NewFileStorage(zap.NewNop())

		service := &service{
			client:          mockClient,
			store:           storage,
			sessionProvider: mockSessionProvider,
		}

		ctx := context.Background()
		fileID := "test-file-id"
		userID := "test-user-id"

		session := authinterfaces.NewSession("token", "key", userID)
		mockSessionProvider.On("GetSession", ctx).Return(session, nil)

		downloadErr := errors.New("download failed")
		mockClient.On("Download", ctx, fileID).Return(nil, nil, downloadErr)

		filepath, err := service.Download(ctx, fileID)

		require.Error(t, err)
		assert.Empty(t, filepath)
		assert.Contains(t, err.Error(), "failed to download file")
		mockClient.AssertExpectations(t)
		mockSessionProvider.AssertExpectations(t)
	})
}

func TestService_List(t *testing.T) {
	t.Parallel()

	t.Run("successful list", func(t *testing.T) {
		mockClient := new(MockClient)
		mockSessionProvider := new(MockSessionProvider)
		storage := filestorage.NewFileStorage(zap.NewNop())

		service := &service{
			client:          mockClient,
			store:           storage,
			sessionProvider: mockSessionProvider,
		}

		ctx := context.Background()
		file1, _ := NewFile("1", "file1.txt", 100, "note1", time.Now())
		file2, _ := NewFile("2", "file2.txt", 200, "note2", time.Now())
		// List returns []File (slice of File structs), so dereference pointers
		expectedFiles := []File{*file1, *file2}

		mockClient.On("List", ctx).Return(expectedFiles, nil)

		files, err := service.List(ctx)

		require.NoError(t, err)
		assert.Equal(t, expectedFiles, files)
		mockClient.AssertExpectations(t)
	})

	t.Run("list error", func(t *testing.T) {
		mockClient := new(MockClient)
		mockSessionProvider := new(MockSessionProvider)
		storage := filestorage.NewFileStorage(zap.NewNop())

		service := &service{
			client:          mockClient,
			store:           storage,
			sessionProvider: mockSessionProvider,
		}

		ctx := context.Background()
		listErr := errors.New("list failed")

		mockClient.On("List", ctx).Return(nil, listErr)

		files, err := service.List(ctx)

		require.Error(t, err)
		assert.Nil(t, files)
		mockClient.AssertExpectations(t)
	})
}
