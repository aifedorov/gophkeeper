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
		name         string
		filePath     string
		notes        string
		setupMock    func(*MockClient, *filestorage.MockStorage)
		wantErr      bool
		errMsg       string
		cleanupFiles bool
	}{
		{
			name:     "successful upload",
			filePath: "test.txt",
			notes:    "test notes",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage) {
				tmpFile := createTempFile(t, "test content")
				// #nosec G304
				ms.EXPECT().
					OpenFile(gomock.Any(), "test.txt").
					Return(os.Open(tmpFile))

				mc.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr:      false,
			cleanupFiles: true,
		},
		{
			name:     "file not found",
			filePath: "/nonexistent/file.txt",
			notes:    "test notes",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage) {
				ms.EXPECT().
					OpenFile(gomock.Any(), "/nonexistent/file.txt").
					Return(nil, errors.New("filestorage: file not found: file does not exist"))
			},
			wantErr: true,
			errMsg:  "file not found",
		},
		{
			name:     "file open error",
			filePath: "/root/readonly.txt",
			notes:    "test notes",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage) {
				ms.EXPECT().
					OpenFile(gomock.Any(), "/root/readonly.txt").
					Return(nil, errors.New("filestorage: file not found: permission denied"))
			},
			wantErr: true,
			errMsg:  "file not found",
		},
		{
			name:     "client upload error",
			filePath: "test.txt",
			notes:    "test notes",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage) {
				tmpFile := createTempFile(t, "test content")
				// #nosec G304
				ms.EXPECT().
					OpenFile(gomock.Any(), "test.txt").
					Return(os.Open(tmpFile))

				mc.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("upload failed")).
					Times(1)
			},
			wantErr:      true,
			errMsg:       "upload failed",
			cleanupFiles: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockClient(ctrl)
			mockStorage := filestorage.NewMockStorage(ctrl)
			tt.setupMock(mockClient, mockStorage)

			service := NewService(mockClient, mockStorage, nil)

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
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*MockClient, *filestorage.MockStorage)
		wantFiles []File
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful list",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage) {
				file1, _ := NewFile("1", "file1.txt", 100, "note1", time.Now())
				file2, _ := NewFile("2", "file2.txt", 200, "note2", time.Now())
				expectedFiles := []File{*file1, *file2}

				mc.EXPECT().
					List(gomock.Any()).
					Return(expectedFiles, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "list error",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage) {
				mc.EXPECT().
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
			mockStorage := filestorage.NewMockStorage(ctrl)
			tt.setupMock(mockClient, mockStorage)

			service := NewService(mockClient, mockStorage, nil)

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
		setupMock func(*MockClient, *filestorage.MockStorage, *mockSessionProvider)
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "successful download",
			fileID: "test-file-id",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage, sp *mockSessionProvider) {
				userID := "test-user-id"
				login := "test-login"
				session := authinterfaces.NewSession("token", "key", userID, login)
				sp.setSession(session)
				sp.setError(nil)

				fileContent := "test file content"
				mockReader := &mockReadCloser{Reader: strings.NewReader(fileContent)}
				mockMeta, _ := NewFileMeta("test_new.txt", int64(len(fileContent)), "test notes")

				mc.EXPECT().
					Download(gomock.Any(), "test-file-id").
					Return(mockReader, mockMeta, nil).
					Times(1)

				ms.EXPECT().
					Upload(gomock.Any(), login, "test_new.txt", mockReader).
					Return("/path/to/test_new.txt", nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:   "session provider error",
			fileID: "test-file-id",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage, sp *mockSessionProvider) {
				sp.setSession(authinterfaces.Session{})
				sp.setError(errors.New("session not found"))
			},
			wantErr: true,
			errMsg:  "failed to get session",
		},
		{
			name:   "client download error",
			fileID: "test-file-id",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage, sp *mockSessionProvider) {
				userID := "test-user-id"
				login := "test-login"
				session := authinterfaces.NewSession("token", "key", userID, login)
				sp.setSession(session)
				sp.setError(nil)

				mc.EXPECT().
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
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage, sp *mockSessionProvider) {
				userID := "test-user-id"
				login := "test-login"
				session := authinterfaces.NewSession("token", "key", userID, login)
				sp.setSession(session)
				sp.setError(nil)

				mockReader := &mockReadCloser{Reader: strings.NewReader("content")}
				mc.EXPECT().
					Download(gomock.Any(), "test-file-id").
					Return(mockReader, nil, nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to download file",
		},
		{
			name:   "storage upload error",
			fileID: "test-file-id",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage, sp *mockSessionProvider) {
				userID := "test-user-id"
				login := "test-login"
				session := authinterfaces.NewSession("token", "key", userID, login)
				sp.setSession(session)
				sp.setError(nil)

				fileContent := "test file content"
				mockReader := &mockReadCloser{Reader: strings.NewReader(fileContent)}
				mockMeta, _ := NewFileMeta("test_new.txt", int64(len(fileContent)), "test notes")

				mc.EXPECT().
					Download(gomock.Any(), "test-file-id").
					Return(mockReader, mockMeta, nil).
					Times(1)

				ms.EXPECT().
					Upload(gomock.Any(), login, "test_new.txt", mockReader).
					Return("", errors.New("storage upload failed")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "storage upload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := NewMockClient(ctrl)
			mockStorage := filestorage.NewMockStorage(ctrl)
			mockSessionProvider := &mockSessionProvider{ctrl: ctrl}
			tt.setupMock(mockClient, mockStorage, mockSessionProvider)

			service := NewService(mockClient, mockStorage, mockSessionProvider)

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
		setupMock func(*MockClient, *filestorage.MockStorage)
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "successful deletion",
			fileID: "test-file-id",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage) {
				mc.EXPECT().
					Delete(gomock.Any(), "test-file-id").
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:   "client delete error",
			fileID: "test-file-id",
			setupMock: func(mc *MockClient, ms *filestorage.MockStorage) {
				mc.EXPECT().
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
			mockStorage := filestorage.NewMockStorage(ctrl)
			tt.setupMock(mockClient, mockStorage)

			service := NewService(mockClient, mockStorage, nil)

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
