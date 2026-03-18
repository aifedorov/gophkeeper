package binary

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_Upload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filePath string
		notes    string
		setup    func(*testSetup, string) string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "successful upload",
			filePath: testFileName,
			notes:    testNotes,
			setup: func(s *testSetup, filePath string) string {
				tmpFile := createTempFile(t, testFileContent)
				s.expectUploadSuccess(filePath, tmpFile)
				return tmpFile
			},
			wantErr: false,
		},
		{
			name:     "file not found",
			filePath: "/nonexistent/file.txt",
			notes:    testNotes,
			setup: func(s *testSetup, filePath string) string {
				s.expectUploadFileNotFound(filePath, errors.New("filestorage: file not found: file does not exist"))
				return ""
			},
			wantErr: true,
			errMsg:  "file not found",
		},
		{
			name:     "file open error",
			filePath: "/root/readonly.txt",
			notes:    testNotes,
			setup: func(s *testSetup, filePath string) string {
				s.expectUploadFileNotFound(filePath, errors.New("filestorage: file not found: permission denied"))
				return ""
			},
			wantErr: true,
			errMsg:  "file not found",
		},
		{
			name:     "client upload error",
			filePath: testFileName,
			notes:    testNotes,
			setup: func(s *testSetup, filePath string) string {
				tmpFile := createTempFile(t, testFileContent)
				s.expectUploadClientError(filePath, tmpFile, errors.New("upload failed"))
				return tmpFile
			},
			wantErr: true,
			errMsg:  "upload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()

			tmpFile := tt.setup(s, tt.filePath)
			if tmpFile != "" {
				defer func() {
					_ = os.Remove(tmpFile)
				}()
			}

			err := s.service.Upload(s.ctx, tt.filePath, tt.notes)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*testSetup)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful list",
			setup: func(s *testSetup) {
				now := time.Now()
				file1, _ := NewFile("1", "file1.txt", 100, "note1", 1, now)
				file2, _ := NewFile("2", "file2.txt", 200, "note2", 1, now)
				expectedFiles := []File{*file1, *file2}
				s.expectListSuccess(expectedFiles)
				s.wantFiles = expectedFiles
			},
			wantErr: false,
		},
		{
			name: "list error",
			setup: func(s *testSetup) {
				s.expectListError(errors.New("list failed"))
			},
			wantErr: true,
			errMsg:  "list failed",
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
			if !tt.wantErr && s.wantFiles != nil {
				assertFilesEqual(t, files, s.wantFiles)
			}
		})
	}
}

func TestService_Download(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fileID   string
		setup    func(*testSetup, string)
		wantErr  bool
		errMsg   string
		wantPath string
	}{
		{
			name:   "successful download",
			fileID: testFileID,
			setup: func(s *testSetup, fileID string) {
				s.expectDownloadSuccess(fileID, testFileName)
			},
			wantErr:  false,
			wantPath: testFilePath,
		},
		{
			name:   "session provider error",
			fileID: testFileID,
			setup: func(s *testSetup, fileID string) {
				s.expectDownloadSessionError(errors.New("session not found"))
			},
			wantErr: true,
			errMsg:  "failed to get session",
		},
		{
			name:   "client download error",
			fileID: testFileID,
			setup: func(s *testSetup, fileID string) {
				s.expectDownloadClientError(fileID, errors.New("download failed"))
			},
			wantErr: true,
			errMsg:  "failed to download file",
		},
		{
			name:   "client download returns nil meta",
			fileID: testFileID,
			setup: func(s *testSetup, fileID string) {
				s.expectDownloadNilMeta(fileID)
			},
			wantErr: true,
			errMsg:  "failed to download file",
		},
		{
			name:   "storage upload error",
			fileID: testFileID,
			setup: func(s *testSetup, fileID string) {
				s.expectDownloadStorageError(fileID, testFileName, errors.New("storage upload failed"))
			},
			wantErr: true,
			errMsg:  "storage upload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s, tt.fileID)

			filepath, err := s.service.Download(s.ctx, tt.fileID)

			assertError(t, err, tt.wantErr, tt.errMsg)
			if !tt.wantErr {
				assert.Equal(t, tt.wantPath, filepath)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fileID   string
		filePath string
		notes    string
		setup    func(*testSetup, string) string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "successful update",
			fileID:   testFileID,
			filePath: testFileName,
			notes:    "updated notes",
			setup: func(s *testSetup, filePath string) string {
				tmpFile := createTempFile(t, testFileContent)
				s.expectUpdateSuccess(filePath, tmpFile)
				return tmpFile
			},
			wantErr: false,
		},
		{
			name:     "file not found",
			fileID:   testFileID,
			filePath: "/nonexistent/file.txt",
			notes:    testNotes,
			setup: func(s *testSetup, filePath string) string {
				s.expectUpdateFileNotFound(filePath, errors.New("filestorage: file not found"))
				return ""
			},
			wantErr: true,
			errMsg:  "file not found",
		},
		{
			name:     "cache get version error",
			fileID:   testFileID,
			filePath: testFileName,
			notes:    testNotes,
			setup: func(s *testSetup, filePath string) string {
				tmpFile := createTempFile(t, testFileContent)
				s.mockStorage.EXPECT().
					OpenFile(gomock.Any(), filePath).
					Return(os.Open(tmpFile)).
					Times(1)
				s.mockCache.getVersionErr = errors.New("version not found")
				return tmpFile
			},
			wantErr: true,
			errMsg:  "failed to get version from cache",
		},
		{
			name:     "client update error",
			fileID:   testFileID,
			filePath: testFileName,
			notes:    testNotes,
			setup: func(s *testSetup, filePath string) string {
				tmpFile := createTempFile(t, testFileContent)
				s.expectUpdateClientError(filePath, tmpFile, errors.New("update failed"))
				return tmpFile
			},
			wantErr: true,
			errMsg:  "failed to update file",
		},
		{
			name:     "cache set version error",
			fileID:   testFileID,
			filePath: testFileName,
			notes:    testNotes,
			setup: func(s *testSetup, filePath string) string {
				tmpFile := createTempFile(t, testFileContent)
				s.mockStorage.EXPECT().
					OpenFile(gomock.Any(), filePath).
					Return(os.Open(tmpFile)).
					Times(1)
				s.mockClient.EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(2), nil).
					Times(1)
				s.mockCache.setVersionErr = errors.New("cache write failed")
				return tmpFile
			},
			wantErr: true,
			errMsg:  "failed to save file version to cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()

			tmpFile := tt.setup(s, tt.filePath)
			if tmpFile != "" {
				defer func() {
					_ = os.Remove(tmpFile)
				}()
			}

			err := s.service.Update(s.ctx, tt.fileID, tt.filePath, tt.notes)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fileID  string
		setup   func(*testSetup, string)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "successful deletion",
			fileID: testFileID,
			setup: func(s *testSetup, fileID string) {
				s.expectDeleteSuccess(fileID)
			},
			wantErr: false,
		},
		{
			name:   "client delete error",
			fileID: testFileID,
			setup: func(s *testSetup, fileID string) {
				s.expectDeleteError(fileID, errors.New("delete failed"))
			},
			wantErr: true,
			errMsg:  "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s, tt.fileID)

			err := s.service.Delete(s.ctx, tt.fileID)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}
