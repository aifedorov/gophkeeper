package binary

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// Test constants
const (
	testFileName = "test-file.txt"
	testFileSize = int64(1024)
	testNotes    = "test notes"
	testUserID   = "test-user-id"
	testFileID   = "test-file-id"
)

var (
	testEncryptionKey    = []byte("test-encryption-key-32-bytes!!")   // 30 bytes, need 32
	testEncryptionKeyStr = "dGVzdC1lbmNyeXB0aW9uLWtleS0zMi1ieXRlcyEh" // base64 encoded
	testFilePath         = "/path/to/test-file.txt"
)

func init() {
	// Ensure key is exactly 32 bytes for AES-256
	key := make([]byte, 32)
	copy(key, "test-encryption-key-32-bytes!!")
	testEncryptionKey = key
	// Update base64 string to match the 32-byte key
	testEncryptionKeyStr = "dGVzdC1lbmNyeXB0aW9uLWtleS0zMi1ieXRlcyEhAAA="
}

type testSetup struct {
	ctrl             *gomock.Controller
	mockRepo         *mocks.MockRepository
	mockFileStore    *mocks.MockFileStorage
	mockCrypto       *mocks.MockCryptoService
	service          Service
	logger           *zap.Logger
	userID           string
	fileID           string
	encryptionKey    []byte
	encryptionKeyStr string
	encryptedPath    []byte
	encryptedSize    []byte
	encryptedNotes   []byte
	fileMetadata     interfaces.FileMetadata
	fileReader       io.Reader
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	return &testSetup{
		ctrl:             ctrl,
		mockRepo:         mocks.NewMockRepository(ctrl),
		mockFileStore:    mocks.NewMockFileStorage(ctrl),
		mockCrypto:       mocks.NewMockCryptoService(ctrl),
		logger:           zap.NewNop(),
		userID:           testUserID,
		fileID:           uuid.New().String(),
		encryptionKey:    testEncryptionKey,
		encryptionKeyStr: testEncryptionKeyStr,
		encryptedPath:    []byte("encrypted-path"),
		encryptedSize:    []byte("encrypted-size"),
		encryptedNotes:   []byte("encrypted-notes"),
		fileMetadata: interfaces.FileMetadata{
			Name:  testFileName,
			Size:  testFileSize,
			Notes: testNotes,
		},
		fileReader: strings.NewReader("test file content"),
	}
}

func (s *testSetup) initService() {
	s.service = NewService(s.mockRepo, s.mockFileStore, s.mockCrypto, s.logger)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
}

func newTestFile(id, name string, size int64, path, notes string) *interfaces.File {
	file, _ := interfaces.NewFile(id, name, size, path, notes, time.Now())
	return file
}

func newTestRepositoryFile(id, name string, encryptedPath, encryptedSize, encryptedNotes []byte) interfaces.RepositoryFile {
	return interfaces.RepositoryFile{
		ID:             id,
		Name:           name,
		EncryptedPath:  encryptedPath,
		EncryptedSize:  encryptedSize,
		EncryptedNotes: encryptedNotes,
		UpdatedAt:      time.Now(),
	}
}

func assertFileFields(t *testing.T, file *interfaces.File, name string, size int64, notes string) {
	t.Helper()
	require.NotNil(t, file)
	assert.Equal(t, name, file.GetName())
	assert.Equal(t, size, file.GetSize())
	assert.Equal(t, notes, file.GetNotes())
}

func assertFileMetadata(t *testing.T, meta interfaces.FileMetadata, name string, size int64, notes string) {
	t.Helper()
	assert.Equal(t, name, meta.Name)
	assert.Equal(t, size, meta.Size)
	assert.Equal(t, notes, meta.Notes)
}
