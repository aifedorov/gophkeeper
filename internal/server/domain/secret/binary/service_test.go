package binary

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Upload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful upload",
			setupMock: func(s *testSetup) {
				s.mockFileStore.EXPECT().
					Upload(gomock.Any(), s.userID, gomock.Any(), gomock.Any()).
					Return(testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(s.encryptedPath, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(testNotes, s.encryptionKey).
					Return(s.encryptedNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(gomock.Any(), s.encryptionKey).
					DoAndReturn(func(text string, key []byte) ([]byte, error) {
						return s.encryptedSize, nil
					}).
					AnyTimes()

				s.mockRepo.EXPECT().
					Create(gomock.Any(), s.userID, gomock.Any()).
					DoAndReturn(func(ctx context.Context, userID string, file interfaces.RepositoryFile) (*interfaces.RepositoryFile, error) {
						file.Version = 1
						return &file, nil
					}).
					Times(1)

				// Decryption for RepositoryToDomain after Create
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return(testNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedSize, s.encryptionKey).
					Return("1024", nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedPath, s.encryptionKey).
					Return(testFilePath, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "invalid encryption key",
			setupMock: func(s *testSetup) {
				// No mocks - validation fails before any calls
			},
			wantErr: true,
			errMsg:  "failed to decode encryption key",
		},
		{
			name: "file store upload failure",
			setupMock: func(s *testSetup) {
				s.mockFileStore.EXPECT().
					Upload(gomock.Any(), s.userID, gomock.Any(), gomock.Any()).
					Return("", errors.New("storage error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to upload file",
		},
		{
			name: "encryption failure for path",
			setupMock: func(s *testSetup) {
				s.mockFileStore.EXPECT().
					Upload(gomock.Any(), s.userID, gomock.Any(), gomock.Any()).
					Return(testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)

				s.mockFileStore.EXPECT().
					Delete(gomock.Any(), s.userID, gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to convert file",
		},
		{
			name: "repository create failure with cleanup",
			setupMock: func(s *testSetup) {
				s.mockFileStore.EXPECT().
					Upload(gomock.Any(), s.userID, gomock.Any(), gomock.Any()).
					Return(testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(s.encryptedPath, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(testNotes, s.encryptionKey).
					Return(s.encryptedNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(gomock.Any(), s.encryptionKey).
					Return(s.encryptedSize, nil).
					Times(1)

				s.mockRepo.EXPECT().
					Create(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, errors.New("db error")).
					Times(1)

				s.mockFileStore.EXPECT().
					Delete(gomock.Any(), s.userID, gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to create file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			if tt.name == "invalid encryption key" {
				setup.encryptionKeyStr = "invalid-base64!!!"
			}

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			result, err := setup.service.Upload(ctx, setup.userID, setup.encryptionKeyStr, setup.fileMetadata, setup.fileReader)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assertFileFields(t, result, testFileName, testFileSize, testNotes)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantCount int
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful list with multiple files",
			setupMock: func(s *testSetup) {
				file1ID := uuid.New().String()
				file2ID := uuid.New().String()

				repoFile1 := newTestRepositoryFile(file1ID, "file1.txt", s.encryptedPath, s.encryptedSize, s.encryptedNotes)
				repoFile2 := newTestRepositoryFile(file2ID, "file2.txt", s.encryptedPath, s.encryptedSize, s.encryptedNotes)

				s.mockRepo.EXPECT().
					List(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryFile{repoFile1, repoFile2}, nil).
					Times(1)

				// Expect decryption for each file (2 files * 3 fields = 6 calls)
				for i := 0; i < 2; i++ {
					s.mockCrypto.EXPECT().
						Decrypt(s.encryptedNotes, s.encryptionKey).
						Return(testNotes, nil).
						Times(1)
					s.mockCrypto.EXPECT().
						Decrypt(s.encryptedSize, s.encryptionKey).
						Return("1024", nil).
						Times(1)
					s.mockCrypto.EXPECT().
						Decrypt(s.encryptedPath, s.encryptionKey).
						Return(testFilePath, nil).
						Times(1)
				}
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "successful list with empty result",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					List(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryFile{}, nil).
					Times(1)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "invalid encryption key",
			setupMock: func(s *testSetup) {
				// No mocks - validation fails before any calls
			},
			wantErr: true,
			errMsg:  "failed to decode encryption key",
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					List(gomock.Any(), s.userID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to list files",
		},
		{
			name: "decryption failure for notes",
			setupMock: func(s *testSetup) {
				repoFile := newTestRepositoryFile(s.fileID, testFileName, s.encryptedPath, s.encryptedSize, s.encryptedNotes)

				s.mockRepo.EXPECT().
					List(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryFile{repoFile}, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to convert file metadata",
		},
		{
			name: "decryption failure for size",
			setupMock: func(s *testSetup) {
				repoFile := newTestRepositoryFile(s.fileID, testFileName, s.encryptedPath, s.encryptedSize, s.encryptedNotes)

				s.mockRepo.EXPECT().
					List(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryFile{repoFile}, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return(testNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedSize, s.encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to convert file metadata",
		},
		{
			name: "decryption failure for path",
			setupMock: func(s *testSetup) {
				repoFile := newTestRepositoryFile(s.fileID, testFileName, s.encryptedPath, s.encryptedSize, s.encryptedNotes)

				s.mockRepo.EXPECT().
					List(gomock.Any(), s.userID).
					Return([]interfaces.RepositoryFile{repoFile}, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return(testNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedSize, s.encryptionKey).
					Return("1024", nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedPath, s.encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to convert file metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			// For invalid encryption key test
			if tt.name == "invalid encryption key" {
				setup.encryptionKeyStr = "invalid-base64!!!"
			}

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			result, err := setup.service.List(ctx, setup.userID, setup.encryptionKeyStr)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}
		})
	}
}

func TestService_Download(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantErr   bool
		wantErrIs error
		errMsg    string
	}{
		{
			name: "successful download",
			setupMock: func(s *testSetup) {
				repoFile := newTestRepositoryFile(s.fileID, testFileName, s.encryptedPath, s.encryptedSize, s.encryptedNotes)

				s.mockRepo.EXPECT().
					Get(gomock.Any(), s.userID, s.fileID).
					Return(&repoFile, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return(testNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedSize, s.encryptionKey).
					Return("1024", nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedPath, s.encryptionKey).
					Return(testFilePath, nil).
					Times(1)

				encryptedReader := io.NopCloser(strings.NewReader("encrypted content"))
				s.mockFileStore.EXPECT().
					Download(gomock.Any(), s.userID, gomock.Any()).
					Return(encryptedReader, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "invalid encryption key",
			setupMock: func(s *testSetup) {
				// No mocks - validation fails before any calls
			},
			wantErr: true,
			errMsg:  "failed to decode encryption key",
		},
		{
			name: "file not found",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					Get(gomock.Any(), s.userID, s.fileID).
					Return(nil, ErrNotFound).
					Times(1)
			},
			wantErr:   true,
			wantErrIs: ErrNotFound,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					Get(gomock.Any(), s.userID, s.fileID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to get file from repository",
		},
		{
			name: "decryption failure",
			setupMock: func(s *testSetup) {
				repoFile := newTestRepositoryFile(s.fileID, testFileName, s.encryptedPath, s.encryptedSize, s.encryptedNotes)

				s.mockRepo.EXPECT().
					Get(gomock.Any(), s.userID, s.fileID).
					Return(&repoFile, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to convert file metadata",
		},
		{
			name: "file store download failure",
			setupMock: func(s *testSetup) {
				repoFile := newTestRepositoryFile(s.fileID, testFileName, s.encryptedPath, s.encryptedSize, s.encryptedNotes)

				s.mockRepo.EXPECT().
					Get(gomock.Any(), s.userID, s.fileID).
					Return(&repoFile, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return(testNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedSize, s.encryptionKey).
					Return("1024", nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedPath, s.encryptionKey).
					Return(testFilePath, nil).
					Times(1)

				s.mockFileStore.EXPECT().
					Download(gomock.Any(), s.userID, gomock.Any()).
					Return(nil, errors.New("file not found")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to open file for reading",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			// For invalid encryption key test
			if tt.name == "invalid encryption key" {
				setup.encryptionKeyStr = "invalid-base64!!!"
			}

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			reader, meta, err := setup.service.Download(ctx, setup.userID, setup.encryptionKeyStr, setup.fileID)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, reader)
			} else if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, reader)
			} else {
				require.NoError(t, err)
				require.NotNil(t, reader)
				assertFileMetadata(t, meta, testFileName, testFileSize, testNotes)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantErr   bool
		wantErrIs error
		errMsg    string
	}{
		{
			name: "successful update",
			setupMock: func(s *testSetup) {
				s.fileMetadata.ID = s.fileID
				s.fileMetadata.Version = 1

				s.mockFileStore.EXPECT().
					BeginUpdate(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return("/tmp/tmpfile", testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(s.encryptedPath, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(testNotes, s.encryptionKey).
					Return(s.encryptedNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(gomock.Any(), s.encryptionKey).
					Return(s.encryptedSize, nil).
					AnyTimes()

				s.mockRepo.EXPECT().
					Update(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					DoAndReturn(func(ctx context.Context, userID, id string, file interfaces.RepositoryFile) (*interfaces.RepositoryFile, error) {
						file.Version = 2
						return &file, nil
					}).
					Times(1)

				// Decryption for RepositoryToDomain
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return(testNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedSize, s.encryptionKey).
					Return("1024", nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedPath, s.encryptionKey).
					Return(testFilePath, nil).
					Times(1)

				s.mockFileStore.EXPECT().
					CommitUpdate(gomock.Any(), s.userID, s.fileID).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "invalid encryption key",
			setupMock: func(s *testSetup) {
				s.fileMetadata.ID = s.fileID
			},
			wantErr: true,
			errMsg:  "failed to wrap reader for encryption",
		},
		{
			name: "begin update failure",
			setupMock: func(s *testSetup) {
				s.fileMetadata.ID = s.fileID

				s.mockFileStore.EXPECT().
					BeginUpdate(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return("", "", errors.New("storage error")).
					Times(1)

				s.mockFileStore.EXPECT().
					AbortUpdate(gomock.Any(), "").
					Return(nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to update file",
		},
		{
			name: "file to repository conversion failure",
			setupMock: func(s *testSetup) {
				s.fileMetadata.ID = s.fileID

				s.mockFileStore.EXPECT().
					BeginUpdate(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return("/tmp/tmpfile", testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(nil, errors.New("encryption error")).
					Times(1)

				s.mockFileStore.EXPECT().
					AbortUpdate(gomock.Any(), "/tmp/tmpfile").
					Return(nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to convert file",
		},
		{
			name: "file not found",
			setupMock: func(s *testSetup) {
				s.fileMetadata.ID = s.fileID

				s.mockFileStore.EXPECT().
					BeginUpdate(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return("/tmp/tmpfile", testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(s.encryptedPath, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(testNotes, s.encryptionKey).
					Return(s.encryptedNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(gomock.Any(), s.encryptionKey).
					Return(s.encryptedSize, nil).
					AnyTimes()

				s.mockRepo.EXPECT().
					Update(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return(nil, ErrNotFound).
					Times(1)

				s.mockFileStore.EXPECT().
					AbortUpdate(gomock.Any(), "/tmp/tmpfile").
					Return(nil).
					Times(1)
			},
			wantErr:   true,
			wantErrIs: ErrNotFound,
		},
		{
			name: "version conflict",
			setupMock: func(s *testSetup) {
				s.fileMetadata.ID = s.fileID

				s.mockFileStore.EXPECT().
					BeginUpdate(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return("/tmp/tmpfile", testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(s.encryptedPath, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(testNotes, s.encryptionKey).
					Return(s.encryptedNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(gomock.Any(), s.encryptionKey).
					Return(s.encryptedSize, nil).
					AnyTimes()

				s.mockRepo.EXPECT().
					Update(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return(nil, ErrVersionConflict).
					Times(1)

				s.mockFileStore.EXPECT().
					AbortUpdate(gomock.Any(), "/tmp/tmpfile").
					Return(nil).
					Times(1)
			},
			wantErr:   true,
			wantErrIs: ErrVersionConflict,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.fileMetadata.ID = s.fileID

				s.mockFileStore.EXPECT().
					BeginUpdate(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return("/tmp/tmpfile", testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(s.encryptedPath, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(testNotes, s.encryptionKey).
					Return(s.encryptedNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(gomock.Any(), s.encryptionKey).
					Return(s.encryptedSize, nil).
					AnyTimes()

				s.mockRepo.EXPECT().
					Update(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return(nil, errors.New("db error")).
					Times(1)

				s.mockFileStore.EXPECT().
					AbortUpdate(gomock.Any(), "/tmp/tmpfile").
					Return(nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to update file",
		},
		{
			name: "repository to domain conversion failure",
			setupMock: func(s *testSetup) {
				s.fileMetadata.ID = s.fileID

				s.mockFileStore.EXPECT().
					BeginUpdate(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return("/tmp/tmpfile", testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(s.encryptedPath, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(testNotes, s.encryptionKey).
					Return(s.encryptedNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(gomock.Any(), s.encryptionKey).
					Return(s.encryptedSize, nil).
					AnyTimes()

				s.mockRepo.EXPECT().
					Update(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					DoAndReturn(func(ctx context.Context, userID, id string, file interfaces.RepositoryFile) (*interfaces.RepositoryFile, error) {
						file.Version = 2
						return &file, nil
					}).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return("", errors.New("decryption error")).
					Times(1)

				s.mockFileStore.EXPECT().
					AbortUpdate(gomock.Any(), "/tmp/tmpfile").
					Return(nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to convert repository file to domain",
		},
		{
			name: "commit update failure",
			setupMock: func(s *testSetup) {
				s.fileMetadata.ID = s.fileID

				s.mockFileStore.EXPECT().
					BeginUpdate(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					Return("/tmp/tmpfile", testFilePath, nil).
					Times(1)

				s.mockCrypto.EXPECT().
					Encrypt(testFilePath, s.encryptionKey).
					Return(s.encryptedPath, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(testNotes, s.encryptionKey).
					Return(s.encryptedNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Encrypt(gomock.Any(), s.encryptionKey).
					Return(s.encryptedSize, nil).
					AnyTimes()

				s.mockRepo.EXPECT().
					Update(gomock.Any(), s.userID, s.fileID, gomock.Any()).
					DoAndReturn(func(ctx context.Context, userID, id string, file interfaces.RepositoryFile) (*interfaces.RepositoryFile, error) {
						file.Version = 2
						return &file, nil
					}).
					Times(1)

				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedNotes, s.encryptionKey).
					Return(testNotes, nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedSize, s.encryptionKey).
					Return("1024", nil).
					Times(1)
				s.mockCrypto.EXPECT().
					Decrypt(s.encryptedPath, s.encryptionKey).
					Return(testFilePath, nil).
					Times(1)

				s.mockFileStore.EXPECT().
					CommitUpdate(gomock.Any(), s.userID, s.fileID).
					Return(errors.New("commit error")).
					Times(1)

				s.mockFileStore.EXPECT().
					AbortUpdate(gomock.Any(), "/tmp/tmpfile").
					Return(nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to commit file update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			if tt.name == "invalid encryption key" {
				setup.encryptionKeyStr = "invalid-base64!!!"
			}

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			result, err := setup.service.Update(ctx, setup.userID, setup.encryptionKeyStr, setup.fileMetadata, setup.fileReader)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
				assert.Nil(t, result)
			} else if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assertFileFields(t, result, testFileName, testFileSize, testNotes)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(*testSetup)
		wantErr   bool
		wantErrIs error
		errMsg    string
	}{
		{
			name: "successful deletion",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					Delete(gomock.Any(), s.userID, s.fileID).
					Return(nil).
					Times(1)

				s.mockFileStore.EXPECT().
					Delete(gomock.Any(), s.userID, s.fileID).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "file not found",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					Delete(gomock.Any(), s.userID, s.fileID).
					Return(ErrNotFound).
					Times(1)
			},
			wantErr:   true,
			wantErrIs: ErrNotFound,
		},
		{
			name: "repository error",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					Delete(gomock.Any(), s.userID, s.fileID).
					Return(errors.New("db error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to delete file from repository",
		},
		{
			name: "file store delete failure (non-fatal)",
			setupMock: func(s *testSetup) {
				s.mockRepo.EXPECT().
					Delete(gomock.Any(), s.userID, s.fileID).
					Return(nil).
					Times(1)

				s.mockFileStore.EXPECT().
					Delete(gomock.Any(), s.userID, s.fileID).
					Return(errors.New("file system error")).
					Times(1)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newTestSetup(t)
			defer setup.cleanup()

			tt.setupMock(setup)
			setup.initService()

			ctx := context.Background()
			err := setup.service.Delete(ctx, setup.userID, setup.fileID)

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
			} else if tt.wantErr {
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
