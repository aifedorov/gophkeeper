package text

import (
	"errors"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	s := newTestSetup(t)
	defer s.cleanup()
	service := NewService(s.mockBinarySrv, s.mockStorage)

	require.NotNil(t, service)
}

func TestService_CreateFromContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		title   string
		notes   string
		setup   func(*testSetup)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "successful creation",
			content: testContent,
			title:   testTitle,
			notes:   testNotes,
			setup: func(s *testSetup) {
				s.expectCreateFromContentSuccess(testTitle)
			},
			wantErr: false,
		},
		{
			name:    "successful creation without notes",
			content: "Some content",
			title:   "note-without-notes",
			notes:   "",
			setup: func(s *testSetup) {
				s.expectCreateFromContentSuccess("note-without-notes")
			},
			wantErr: false,
		},
		{
			name:    "storage error",
			content: testContent,
			title:   testTitle,
			notes:   testNotes,
			setup: func(s *testSetup) {
				s.expectCreateFromContentStorageError(testTitle, errors.New("storage error"))
			},
			wantErr: true,
			errMsg:  "failed to create file from content",
		},
		{
			name:    "upload error",
			content: testContent,
			title:   testTitle,
			notes:   testNotes,
			setup: func(s *testSetup) {
				s.expectCreateFromContentUploadError(testTitle, errors.New("upload error"))
			},
			wantErr: true,
			errMsg:  "upload error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			err := s.service.CreateFromContent(s.ctx, tt.content, tt.title, tt.notes)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_CreateFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filePath string
		notes    string
		setup    func(*testSetup)
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "successful creation",
			filePath: testFilePath,
			notes:    testNotes,
			setup: func(s *testSetup) {
				s.expectCreateFromFileSuccess()
			},
			wantErr: false,
		},
		{
			name:     "upload error",
			filePath: testFilePath,
			notes:    testNotes,
			setup: func(s *testSetup) {
				s.expectCreateFromFileError(errors.New("upload error"))
			},
			wantErr: true,
			errMsg:  "upload error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			err := s.service.CreateFromFile(s.ctx, tt.filePath, tt.notes)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_View(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          string
		setup       func(*testSetup)
		wantContent string
		wantErr     bool
		errMsg      string
	}{
		{
			name: "successful view",
			id:   testID,
			setup: func(s *testSetup) {
				s.expectViewSuccess(testID, "/tmp/downloaded.txt", testContent)
			},
			wantContent: testContent,
			wantErr:     false,
		},
		{
			name: "download error",
			id:   testID,
			setup: func(s *testSetup) {
				s.expectViewDownloadError(testID, errors.New("download error"))
			},
			wantErr: true,
			errMsg:  "failed to download file",
		},
		{
			name: "read content error",
			id:   testID,
			setup: func(s *testSetup) {
				s.expectViewReadError(testID, "/tmp/downloaded.txt", errors.New("read error"))
			},
			wantErr: true,
			errMsg:  "failed to read file content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			content, err := s.service.View(s.ctx, tt.id)

			assertError(t, err, tt.wantErr, tt.errMsg)
			if !tt.wantErr {
				assert.Equal(t, tt.wantContent, content)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(*testSetup)
		wantCount int
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful list with files",
			setup: func(s *testSetup) {
				files := []binary.File{}
				s.expectListSuccess(files)
				s.wantFiles = files
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "list error",
			setup: func(s *testSetup) {
				s.expectListError(errors.New("list error"))
			},
			wantErr: true,
			errMsg:  "list error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			files, err := s.service.List(s.ctx)

			assertError(t, err, tt.wantErr, tt.errMsg)
			if !tt.wantErr {
				assert.Len(t, files, tt.wantCount)
			}
		})
	}
}

func TestService_UpdateFromContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		content string
		title   string
		notes   string
		setup   func(*testSetup)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "successful update",
			id:      testID,
			content: "updated content",
			title:   "updated-title",
			notes:   "updated notes",
			setup: func(s *testSetup) {
				s.expectUpdateFromContentSuccess(testID, "updated-title")
			},
			wantErr: false,
		},
		{
			name:    "storage error",
			id:      testID,
			content: testContent,
			title:   testTitle,
			notes:   testNotes,
			setup: func(s *testSetup) {
				s.expectUpdateFromContentStorageError(testTitle, errors.New("storage error"))
			},
			wantErr: true,
			errMsg:  "failed to create file from content",
		},
		{
			name:    "update error",
			id:      testID,
			content: testContent,
			title:   testTitle,
			notes:   testNotes,
			setup: func(s *testSetup) {
				s.expectUpdateFromContentUpdateError(testID, testTitle, errors.New("update error"))
			},
			wantErr: true,
			errMsg:  "update error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			err := s.service.UpdateFromContent(s.ctx, tt.id, tt.content, tt.title, tt.notes)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_UpdateFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id       string
		filePath string
		notes    string
		setup    func(*testSetup)
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "successful update",
			id:       testID,
			filePath: testFilePath,
			notes:    testNotes,
			setup: func(s *testSetup) {
				s.expectUpdateFromFileSuccess(testID)
			},
			wantErr: false,
		},
		{
			name:     "update error",
			id:       testID,
			filePath: testFilePath,
			notes:    testNotes,
			setup: func(s *testSetup) {
				s.expectUpdateFromFileError(testID, errors.New("update error"))
			},
			wantErr: true,
			errMsg:  "update error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			err := s.service.UpdateFromFile(s.ctx, tt.id, tt.filePath, tt.notes)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_Download(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		id           string
		setup        func(*testSetup)
		wantFilePath string
		wantErr      bool
		errMsg       string
	}{
		{
			name: "successful download",
			id:   testID,
			setup: func(s *testSetup) {
				s.expectDownloadSuccess(testID, "/downloads/file.txt")
			},
			wantFilePath: "/downloads/file.txt",
			wantErr:      false,
		},
		{
			name: "download error",
			id:   testID,
			setup: func(s *testSetup) {
				s.expectDownloadError(testID, errors.New("download error"))
			},
			wantErr: true,
			errMsg:  "download error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			filePath, err := s.service.Download(s.ctx, tt.id)

			assertError(t, err, tt.wantErr, tt.errMsg)
			if !tt.wantErr {
				assert.Equal(t, tt.wantFilePath, filePath)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		setup   func(*testSetup)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful deletion",
			id:   testID,
			setup: func(s *testSetup) {
				s.expectDeleteSuccess(testID)
			},
			wantErr: false,
		},
		{
			name: "delete error",
			id:   testID,
			setup: func(s *testSetup) {
				s.expectDeleteError(testID, errors.New("delete error"))
			},
			wantErr: true,
			errMsg:  "delete error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			err := s.service.Delete(s.ctx, tt.id)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}
