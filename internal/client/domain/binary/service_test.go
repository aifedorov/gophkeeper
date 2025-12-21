package binary

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	authinterfaces "github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/pkg/filestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// mockReadCloser is a simple implementation of io.ReadCloser for testing
type mockReadCloser struct {
	io.Reader
	closed bool
}

func (m *mockReadCloser) Close() error {
	m.closed = true
	return nil
}

func TestService_Upload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		filePath  string
		notes     string
		setupMock func(*MockClient)
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "successful upload",
			filePath: createTempFile(t, "test content"),
			notes:    "test notes",
			setupMock: func(m *MockClient) {
				m.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:     "file not found",
			filePath: "/nonexistent/file.txt",
			notes:    "test notes",
			setupMock: func(m *MockClient) {
				// No expectation - file open fails before client call
			},
			wantErr: true,
			errMsg:  "file not found",
		},
		{
			name:     "file open error",
			filePath: "/root/readonly.txt",
			notes:    "test notes",
			setupMock: func(m *MockClient) {
				// No expectation - file open fails before client call
			},
			wantErr: true,
			errMsg:  "file not found",
		},
		{
			name:     "client upload error",
			filePath: createTempFile(t, "test content"),
			notes:    "test notes",
			setupMock: func(m *MockClient) {
				m.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("upload failed")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "upload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockClient(ctrl)
			tt.setupMock(mockClient)

			storage := filestorage.NewFileStorage(zap.NewNop())
			service := NewService(mockClient, storage, nil)

			ctx := context.Background()
			err := service.Upload(ctx, tt.filePath, tt.notes)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
			}

			if strings.HasPrefix(tt.filePath, os.TempDir()) {
				_ = os.Remove(tt.filePath)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*MockClient)
		wantFiles []File
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful list",
			setupMock: func(m *MockClient) {
				file1, _ := NewFile("1", "file1.txt", 100, "note1", time.Now())
				file2, _ := NewFile("2", "file2.txt", 200, "note2", time.Now())
				expectedFiles := []File{*file1, *file2}

				m.EXPECT().
					List(gomock.Any()).
					Return(expectedFiles, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "list error",
			setupMock: func(m *MockClient) {
				m.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("list failed")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "list failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockClient(ctrl)
			tt.setupMock(mockClient)

			storage := filestorage.NewFileStorage(zap.NewNop())
			service := NewService(mockClient, storage, nil)

			ctx := context.Background()
			files, err := service.List(ctx)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, files)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, files)
			}
		})
	}
}

func TestService_Download(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fileID    string
		setupMock func(*MockClient, *mockSessionProvider)
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "successful download",
			fileID: "test-file-id",
			setupMock: func(m *MockClient, sp *mockSessionProvider) {
				userID := "test-user-id"
				session := authinterfaces.NewSession("token", "key", userID)
				sp.setSession(session)
				sp.setError(nil)

				fileContent := "test file content"
				mockReader := &mockReadCloser{Reader: strings.NewReader(fileContent)}
				mockMeta, _ := NewFileMeta("test.txt", int64(len(fileContent)), "test notes")

				m.EXPECT().
					Download(gomock.Any(), "test-file-id").
					Return(mockReader, mockMeta, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:   "session provider error",
			fileID: "test-file-id",
			setupMock: func(m *MockClient, sp *mockSessionProvider) {
				sp.setSession(authinterfaces.Session{})
				sp.setError(errors.New("session not found"))
			},
			wantErr: true,
			errMsg:  "failed to get session",
		},
		{
			name:   "client download error",
			fileID: "test-file-id",
			setupMock: func(m *MockClient, sp *mockSessionProvider) {
				userID := "test-user-id"
				session := authinterfaces.NewSession("token", "key", userID)
				sp.setSession(session)
				sp.setError(nil)

				m.EXPECT().
					Download(gomock.Any(), "test-file-id").
					Return(nil, nil, errors.New("download failed")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to download file",
		},
		{
			name:   "client download returns nil meta",
			fileID: "test-file-id",
			setupMock: func(m *MockClient, sp *mockSessionProvider) {
				userID := "test-user-id"
				session := authinterfaces.NewSession("token", "key", userID)
				sp.setSession(session)
				sp.setError(nil)

				mockReader := &mockReadCloser{Reader: strings.NewReader("content")}
				m.EXPECT().
					Download(gomock.Any(), "test-file-id").
					Return(mockReader, nil, nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to download file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockClient(ctrl)
			mockSessionProvider := &mockSessionProvider{ctrl: ctrl}
			tt.setupMock(mockClient, mockSessionProvider)

			storage := filestorage.NewFileStorage(zap.NewNop())
			service := NewService(mockClient, storage, mockSessionProvider)

			ctx := context.Background()
			filepath, err := service.Download(ctx, tt.fileID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Empty(t, filepath)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, filepath)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fileID    string
		setupMock func(*MockClient)
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "successful deletion",
			fileID: "test-file-id",
			setupMock: func(m *MockClient) {
				m.EXPECT().
					Delete(gomock.Any(), "test-file-id").
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:   "client delete error",
			fileID: "test-file-id",
			setupMock: func(m *MockClient) {
				m.EXPECT().
					Delete(gomock.Any(), "test-file-id").
					Return(errors.New("delete failed")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockClient(ctrl)
			tt.setupMock(mockClient)

			storage := filestorage.NewFileStorage(zap.NewNop())
			service := NewService(mockClient, storage, nil)

			ctx := context.Background()
			err := service.Delete(ctx, tt.fileID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// mockSessionProvider is a simple mock for SessionProvider using gomock pattern
type mockSessionProvider struct {
	ctrl    *gomock.Controller
	session authinterfaces.Session
	err     error
	expects []func()
}

func (m *mockSessionProvider) GetSession(ctx context.Context) (authinterfaces.Session, error) {
	return m.session, m.err
}

func (m *mockSessionProvider) setSession(session authinterfaces.Session) {
	m.session = session
}

func (m *mockSessionProvider) setError(err error) {
	m.err = err
}

// Helper function to create a temporary file for testing
func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	require.NoError(t, err)

	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)

	err = tmpfile.Close()
	require.NoError(t, err)

	return tmpfile.Name()
}
