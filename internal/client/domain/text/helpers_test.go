package text

import (
	"context"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/aifedorov/gophkeeper/internal/client/domain/text/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testID       = "test-id-123"
	testTitle    = "test-note"
	testContent  = "This is test content"
	testNotes    = "test notes"
	testFilePath = "/tmp/test.txt"
)

type testSetup struct {
	ctrl          *gomock.Controller
	mockBinarySrv *MockService
	mockStorage   *interfaces.MockStorage
	service       Service
	ctx           context.Context
	wantFiles     []binary.File
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	return &testSetup{
		ctrl:          ctrl,
		mockBinarySrv: NewMockService(ctrl),
		mockStorage:   interfaces.NewMockStorage(ctrl),
		ctx:           context.Background(),
	}
}

func (s *testSetup) initService() {
	s.service = NewService(s.mockBinarySrv, s.mockStorage)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
}

func (s *testSetup) expectCreateFromContentSuccess(title string) {
	s.mockStorage.EXPECT().
		Upload(gomock.Any(), ".tmp", title+".txt", gomock.Any()).
		Return("/tmp/"+title+".txt", nil).
		Times(1)
	s.mockBinarySrv.EXPECT().
		Upload(gomock.Any(), "/tmp/"+title+".txt", gomock.Any()).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectCreateFromContentStorageError(title string, err error) {
	s.mockStorage.EXPECT().
		Upload(gomock.Any(), ".tmp", title+".txt", gomock.Any()).
		Return("", err).
		Times(1)
}

func (s *testSetup) expectCreateFromContentUploadError(title string, err error) {
	s.mockStorage.EXPECT().
		Upload(gomock.Any(), ".tmp", title+".txt", gomock.Any()).
		Return("/tmp/"+title+".txt", nil).
		Times(1)
	s.mockBinarySrv.EXPECT().
		Upload(gomock.Any(), "/tmp/"+title+".txt", gomock.Any()).
		Return(err).
		Times(1)
}

func (s *testSetup) expectCreateFromFileSuccess() {
	s.mockBinarySrv.EXPECT().
		Upload(gomock.Any(), testFilePath, gomock.Any()).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectCreateFromFileError(err error) {
	s.mockBinarySrv.EXPECT().
		Upload(gomock.Any(), testFilePath, gomock.Any()).
		Return(err).
		Times(1)
}

func (s *testSetup) expectViewSuccess(id, path, content string) {
	s.mockBinarySrv.EXPECT().
		Download(gomock.Any(), id).
		Return(path, nil).
		Times(1)
	s.mockStorage.EXPECT().
		ReadContent(gomock.Any(), path, int64(MaxViewSize)).
		Return(content, nil).
		Times(1)
}

func (s *testSetup) expectViewDownloadError(id string, err error) {
	s.mockBinarySrv.EXPECT().
		Download(gomock.Any(), id).
		Return("", err).
		Times(1)
}

func (s *testSetup) expectViewReadError(id, path string, err error) {
	s.mockBinarySrv.EXPECT().
		Download(gomock.Any(), id).
		Return(path, nil).
		Times(1)
	s.mockStorage.EXPECT().
		ReadContent(gomock.Any(), path, int64(MaxViewSize)).
		Return("", err).
		Times(1)
}

func (s *testSetup) expectListSuccess(files []binary.File) {
	s.mockBinarySrv.EXPECT().
		List(gomock.Any()).
		Return(files, nil).
		Times(1)
}

func (s *testSetup) expectListError(err error) {
	s.mockBinarySrv.EXPECT().
		List(gomock.Any()).
		Return(nil, err).
		Times(1)
}

func (s *testSetup) expectUpdateFromContentSuccess(id, title string) {
	s.mockStorage.EXPECT().
		Upload(gomock.Any(), ".tmp", title+".txt", gomock.Any()).
		Return("/tmp/"+title+".txt", nil).
		Times(1)
	s.mockBinarySrv.EXPECT().
		Update(gomock.Any(), id, "/tmp/"+title+".txt", gomock.Any()).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectUpdateFromContentStorageError(title string, err error) {
	s.mockStorage.EXPECT().
		Upload(gomock.Any(), ".tmp", title+".txt", gomock.Any()).
		Return("", err).
		Times(1)
}

func (s *testSetup) expectUpdateFromContentUpdateError(id, title string, err error) {
	s.mockStorage.EXPECT().
		Upload(gomock.Any(), ".tmp", title+".txt", gomock.Any()).
		Return("/tmp/"+title+".txt", nil).
		Times(1)
	s.mockBinarySrv.EXPECT().
		Update(gomock.Any(), id, "/tmp/"+title+".txt", gomock.Any()).
		Return(err).
		Times(1)
}

func (s *testSetup) expectUpdateFromFileSuccess(id string) {
	s.mockBinarySrv.EXPECT().
		Update(gomock.Any(), id, testFilePath, gomock.Any()).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectUpdateFromFileError(id string, err error) {
	s.mockBinarySrv.EXPECT().
		Update(gomock.Any(), id, testFilePath, gomock.Any()).
		Return(err).
		Times(1)
}

func (s *testSetup) expectDownloadSuccess(id, path string) {
	s.mockBinarySrv.EXPECT().
		Download(gomock.Any(), id).
		Return(path, nil).
		Times(1)
}

func (s *testSetup) expectDownloadError(id string, err error) {
	s.mockBinarySrv.EXPECT().
		Download(gomock.Any(), id).
		Return("", err).
		Times(1)
}

func (s *testSetup) expectDeleteSuccess(id string) {
	s.mockBinarySrv.EXPECT().
		Delete(gomock.Any(), id).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectDeleteError(id string, err error) {
	s.mockBinarySrv.EXPECT().
		Delete(gomock.Any(), id).
		Return(err).
		Times(1)
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
