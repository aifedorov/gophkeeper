package binary

import (
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMetadataToFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		metadata interfaces.FileMetadata
		wantErr  bool
	}{
		{
			name: "successful conversion",
			metadata: interfaces.FileMetadata{
				Name:  "test.txt",
				Size:  1024,
				Notes: "test notes",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			metadata: interfaces.FileMetadata{
				Name:  "",
				Size:  1024,
				Notes: "test notes",
			},
			wantErr: true,
		},
		{
			name: "zero size",
			metadata: interfaces.FileMetadata{
				Name:  "test.txt",
				Size:  0,
				Notes: "test notes",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := MetadataToFile(tt.metadata)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.metadata.Name, result.GetName())
				assert.Equal(t, tt.metadata.Size, result.GetSize())
				assert.Equal(t, tt.metadata.Notes, result.GetNotes())
			}
		})
	}
}

func TestFileToRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		file      *interfaces.File
		setupMock func(*mocks.MockCryptoService)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful conversion",
			file: newTestFile("test-id", "test.txt", 1024, "/path/to/file", "notes"),
			setupMock: func(m *mocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("/path/to/file", gomock.Any()).
					Return([]byte("encrypted-path"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("notes", gomock.Any()).
					Return([]byte("encrypted-notes"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("1024", gomock.Any()).
					Return([]byte("encrypted-size"), nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "nil file",
			file: nil,
			setupMock: func(m *mocks.MockCryptoService) {
				// No expectations
			},
			wantErr: true,
			errMsg:  "file is nil",
		},
		{
			name: "encryption failure for path",
			file: newTestFile("test-id", "test.txt", 1024, "/path/to/file", "notes"),
			setupMock: func(m *mocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("/path/to/file", gomock.Any()).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to encrypt path",
		},
		{
			name: "encryption failure for notes",
			file: newTestFile("test-id", "test.txt", 1024, "/path/to/file", "notes"),
			setupMock: func(m *mocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("/path/to/file", gomock.Any()).
					Return([]byte("encrypted-path"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("notes", gomock.Any()).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to encrypt notes",
		},
		{
			name: "encryption failure for size",
			file: newTestFile("test-id", "test.txt", 1024, "/path/to/file", "notes"),
			setupMock: func(m *mocks.MockCryptoService) {
				m.EXPECT().
					Encrypt("/path/to/file", gomock.Any()).
					Return([]byte("encrypted-path"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("notes", gomock.Any()).
					Return([]byte("encrypted-notes"), nil).
					Times(1)
				m.EXPECT().
					Encrypt("1024", gomock.Any()).
					Return(nil, errors.New("encryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to encrypt size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCrypto := mocks.NewMockCryptoService(ctrl)
			tt.setupMock(mockCrypto)

			key := []byte("test-key-32-bytes-long-for-aes!!")
			result, err := FileToRepository(mockCrypto, key, tt.file)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.file.GetID(), result.ID)
				assert.Equal(t, tt.file.GetName(), result.Name)
			}
		})
	}
}

func TestFileToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		repoFile  interfaces.RepositoryFile
		setupMock func(*mocks.MockCryptoService)
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "successful conversion",
			repoFile: newTestRepositoryFile("test-id", "test.txt", []byte("encrypted-path"), []byte("encrypted-size"), []byte("encrypted-notes")),
			setupMock: func(m *mocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-notes"), gomock.Any()).
					Return("notes", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-size"), gomock.Any()).
					Return("1024", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-path"), gomock.Any()).
					Return("/path/to/file", nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:     "decryption failure for notes",
			repoFile: newTestRepositoryFile("test-id", "test.txt", []byte("encrypted-path"), []byte("encrypted-size"), []byte("encrypted-notes")),
			setupMock: func(m *mocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-notes"), gomock.Any()).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to decrypt notes",
		},
		{
			name:     "decryption failure for size",
			repoFile: newTestRepositoryFile("test-id", "test.txt", []byte("encrypted-path"), []byte("encrypted-size"), []byte("encrypted-notes")),
			setupMock: func(m *mocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-notes"), gomock.Any()).
					Return("notes", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-size"), gomock.Any()).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to decrypt size",
		},
		{
			name:     "invalid size format",
			repoFile: newTestRepositoryFile("test-id", "test.txt", []byte("encrypted-path"), []byte("encrypted-size"), []byte("encrypted-notes")),
			setupMock: func(m *mocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-notes"), gomock.Any()).
					Return("notes", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-size"), gomock.Any()).
					Return("not-a-number", nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to convert size",
		},
		{
			name:     "decryption failure for path",
			repoFile: newTestRepositoryFile("test-id", "test.txt", []byte("encrypted-path"), []byte("encrypted-size"), []byte("encrypted-notes")),
			setupMock: func(m *mocks.MockCryptoService) {
				m.EXPECT().
					Decrypt([]byte("encrypted-notes"), gomock.Any()).
					Return("notes", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-size"), gomock.Any()).
					Return("1024", nil).
					Times(1)
				m.EXPECT().
					Decrypt([]byte("encrypted-path"), gomock.Any()).
					Return("", errors.New("decryption error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to decrypt path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCrypto := mocks.NewMockCryptoService(ctrl)
			tt.setupMock(mockCrypto)

			key := []byte("test-key-32-bytes-long-for-aes!!")
			result, err := FileToDomain(mockCrypto, key, tt.repoFile)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.repoFile.ID, result.GetID())
				assert.Equal(t, tt.repoFile.Name, result.GetName())
			}
		})
	}
}

func TestFileToMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     *interfaces.File
		wantErr  bool
		errMsg   string
		wantMeta interfaces.FileMetadata
	}{
		{
			name:    "successful conversion",
			file:    newTestFile("test-id", "test.txt", 1024, "/path/to/file", "notes"),
			wantErr: false,
			wantMeta: interfaces.FileMetadata{
				Name:  "test.txt",
				Size:  1024,
				Notes: "notes",
			},
		},
		{
			name:     "nil file",
			file:     nil,
			wantErr:  true,
			errMsg:   "file is nil",
			wantMeta: interfaces.FileMetadata{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := FileToMetadata(tt.file)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantMeta.Name, result.Name)
				assert.Equal(t, tt.wantMeta.Size, result.Size)
				assert.Equal(t, tt.wantMeta.Notes, result.Notes)
			}
		})
	}
}
