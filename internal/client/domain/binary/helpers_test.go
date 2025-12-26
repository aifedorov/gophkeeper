package binary

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aifedorov/gophkeeper/internal/client/domain/shared"
	"github.com/aifedorov/gophkeeper/pkg/filestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testFileID      = "test-file-id"
	testFileName    = "test.txt"
	testFileSize    = int64(100)
	testNotes       = "test notes"
	testUserID      = "test-user-id"
	testLogin       = "test-login"
	testToken       = "test-token"
	testKey         = "test-key"
	testFilePath    = "/path/to/test.txt"
	testFileContent = "test file content"
)

type mockReadCloser struct {
	io.Reader
	closed bool
}

func (m *mockReadCloser) Close() error {
	m.closed = true
	return nil
}

type mockSessionProvider struct {
	ctrl    *gomock.Controller
	session shared.Session
	err     error
}

func (m *mockSessionProvider) GetSession(_ context.Context) (shared.Session, error) {
	return m.session, m.err
}

func (m *mockSessionProvider) setSession(session shared.Session) {
	m.session = session
}

func (m *mockSessionProvider) setError(err error) {
	m.err = err
}

type testSetup struct {
	ctrl                *gomock.Controller
	mockClient          *MockClient
	mockStorage         *filestorage.MockStorage
	mockSessionProvider *mockSessionProvider
	service             Service
	ctx                 context.Context
	testSession         shared.Session
	testFile            *File
	testFileMeta        *FileMeta
	testFileReader      io.ReadCloser
	wantFiles           []File
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	testFile, _ := NewFile(testFileID, testFileName, testFileSize, testNotes, 1, time.Now())
	testFileMeta, _ := NewFileMeta(testFileName, testFileSize, testNotes, 1)
	testSession := shared.NewSession(testToken, testKey, testUserID, testLogin)

	return &testSetup{
		ctrl:                ctrl,
		mockClient:          NewMockClient(ctrl),
		mockStorage:         filestorage.NewMockStorage(ctrl),
		mockSessionProvider: &mockSessionProvider{ctrl: ctrl},
		ctx:                 context.Background(),
		testSession:         testSession,
		testFile:            testFile,
		testFileMeta:        testFileMeta,
		testFileReader:      &mockReadCloser{Reader: strings.NewReader(testFileContent)},
	}
}

type mockCacheStorage struct{}

func (m *mockCacheStorage) SetFileVersion(id string, version int64) error { return nil }
func (m *mockCacheStorage) GetFileVersion(id string) (int64, error)       { return 1, nil }
func (m *mockCacheStorage) DeleteFileVersion(id string) error             { return nil }

func (s *testSetup) initService() {
	s.service = NewService(s.mockClient, s.mockStorage, &mockCacheStorage{}, s.mockSessionProvider)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
}

// #nosec G304
func (s *testSetup) expectUploadSuccess(filePath string, tmpFile string) {
	s.mockStorage.EXPECT().
		OpenFile(gomock.Any(), filePath).
		Return(os.Open(tmpFile)).
		Times(1)
	s.mockClient.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("test-id", int64(1), nil).
		Times(1)
}

// #nosec G304
func (s *testSetup) expectUploadFileNotFound(filePath string, err error) {
	s.mockStorage.EXPECT().
		OpenFile(gomock.Any(), filePath).
		Return(nil, err).
		Times(1)
}

// #nosec G304
func (s *testSetup) expectUploadClientError(filePath string, tmpFile string, err error) {
	s.mockStorage.EXPECT().
		OpenFile(gomock.Any(), filePath).
		Return(os.Open(tmpFile)).
		Times(1)
	s.mockClient.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", int64(0), err).
		Times(1)
}

// #nosec G304
func (s *testSetup) expectListSuccess(files []File) {
	s.mockClient.EXPECT().
		List(gomock.Any()).
		Return(files, nil).
		Times(1)
}

func (s *testSetup) expectListError(err error) {
	s.mockClient.EXPECT().
		List(gomock.Any()).
		Return(nil, err).
		Times(1)
}

func (s *testSetup) expectDownloadSuccess(fileID string, fileName string) {
	s.mockSessionProvider.setSession(s.testSession)
	s.mockSessionProvider.setError(nil)
	s.mockClient.EXPECT().
		Download(gomock.Any(), fileID).
		Return(s.testFileReader, s.testFileMeta, nil).
		Times(1)
	s.mockStorage.EXPECT().
		Upload(gomock.Any(), testLogin, fileName, s.testFileReader).
		Return(testFilePath, nil).
		Times(1)
}

func (s *testSetup) expectDownloadSessionError(err error) {
	s.mockSessionProvider.setSession(shared.Session{})
	s.mockSessionProvider.setError(err)
}

func (s *testSetup) expectDownloadClientError(fileID string, err error) {
	s.mockSessionProvider.setSession(s.testSession)
	s.mockSessionProvider.setError(nil)
	s.mockClient.EXPECT().
		Download(gomock.Any(), fileID).
		Return(nil, nil, err).
		Times(1)
}

func (s *testSetup) expectDownloadNilMeta(fileID string) {
	s.mockSessionProvider.setSession(s.testSession)
	s.mockSessionProvider.setError(nil)
	s.mockClient.EXPECT().
		Download(gomock.Any(), fileID).
		Return(s.testFileReader, nil, nil).
		Times(1)
}

func (s *testSetup) expectDownloadStorageError(fileID string, fileName string, err error) {
	s.mockSessionProvider.setSession(s.testSession)
	s.mockSessionProvider.setError(nil)
	s.mockClient.EXPECT().
		Download(gomock.Any(), fileID).
		Return(s.testFileReader, s.testFileMeta, nil).
		Times(1)
	s.mockStorage.EXPECT().
		Upload(gomock.Any(), testLogin, fileName, s.testFileReader).
		Return("", err).
		Times(1)
}

// expectDeleteSuccess sets up mocks for successful deletion
func (s *testSetup) expectDeleteSuccess(fileID string) {
	s.mockClient.EXPECT().
		Delete(gomock.Any(), fileID).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectDeleteError(fileID string, err error) {
	s.mockClient.EXPECT().
		Delete(gomock.Any(), fileID).
		Return(err).
		Times(1)
}

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

func assertError(t *testing.T, err error, wantErr bool, errMsg string) {
	t.Helper()
	if wantErr {
		require.Error(t, err)
		if errMsg != "" {
			assert.Contains(t, err.Error(), errMsg)
		}
	} else {
		require.NoError(t, err)
	}
}

func assertFilesEqual(t *testing.T, got, want []File) {
	t.Helper()
	assert.Equal(t, want, got)
}
