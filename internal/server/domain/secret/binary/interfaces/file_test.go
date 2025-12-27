package interfaces

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFile(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name       string
		id         string
		fileName   string
		size       int64
		path       string
		notes      string
		version    int64
		uploadedAt time.Time
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "creates file with all fields",
			id:         "file-id-1",
			fileName:   "test.txt",
			size:       1024,
			path:       "/storage/files/test.txt",
			notes:      "test notes",
			version:    1,
			uploadedAt: now,
			wantErr:    false,
		},
		{
			name:       "creates file without notes",
			id:         "file-id-2",
			fileName:   "document.pdf",
			size:       2048,
			path:       "/storage/files/document.pdf",
			notes:      "",
			version:    2,
			uploadedAt: now,
			wantErr:    false,
		},
		{
			name:       "creates file without path",
			id:         "file-id-3",
			fileName:   "image.png",
			size:       4096,
			path:       "",
			notes:      "image",
			version:    1,
			uploadedAt: now,
			wantErr:    false,
		},
		{
			name:       "empty id",
			id:         "",
			fileName:   "test.txt",
			size:       1024,
			path:       "/path",
			notes:      "notes",
			version:    1,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "file id is required",
		},
		{
			name:       "empty file name",
			id:         "file-id-4",
			fileName:   "",
			size:       1024,
			path:       "/path",
			notes:      "notes",
			version:    1,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "file name is required",
		},
		{
			name:       "zero size",
			id:         "file-id-5",
			fileName:   "test.txt",
			size:       0,
			path:       "/path",
			notes:      "notes",
			version:    1,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "file size is required",
		},
		{
			name:       "size exceeds maximum",
			id:         "file-id-6",
			fileName:   "large.bin",
			size:       11 * 1024 * 1024 * 1024, // 11GB
			path:       "/path",
			notes:      "notes",
			version:    1,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "file size exceeds maximum allowed size",
		},
		{
			name:       "zero version",
			id:         "file-id-7",
			fileName:   "test.txt",
			size:       1024,
			path:       "/path",
			notes:      "notes",
			version:    0,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "invalid file version",
		},
		{
			name:       "negative version",
			id:         "file-id-8",
			fileName:   "test.txt",
			size:       1024,
			path:       "/path",
			notes:      "notes",
			version:    -1,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "invalid file version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := NewFile(tt.id, tt.fileName, tt.size, tt.path, tt.notes, tt.version, tt.uploadedAt)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, file)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, file)
			assert.Equal(t, tt.id, file.GetID())
			assert.Equal(t, tt.fileName, file.GetName())
			assert.Equal(t, tt.size, file.GetSize())
			assert.Equal(t, tt.path, file.GetPath())
			assert.Equal(t, tt.notes, file.GetNotes())
			assert.Equal(t, tt.version, file.GetVersion())
			assert.Equal(t, tt.uploadedAt, file.GetUploadedAt())
		})
	}
}

func TestFile_SetPath(t *testing.T) {
	t.Parallel()

	now := time.Now()
	file, err := NewFile("test-id", "test.txt", 1024, "", "notes", 1, now)
	require.NoError(t, err)

	assert.Equal(t, "", file.GetPath())

	file.SetPath("/new/path/test.txt")
	assert.Equal(t, "/new/path/test.txt", file.GetPath())
}
