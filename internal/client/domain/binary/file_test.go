package binary

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
			notes:      "",
			version:    2,
			uploadedAt: now,
			wantErr:    false,
		},
		{
			name:       "empty file name",
			id:         "file-id-3",
			fileName:   "",
			size:       1024,
			notes:      "notes",
			version:    1,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "file name is required",
		},
		{
			name:       "zero size",
			id:         "file-id-4",
			fileName:   "test.txt",
			size:       0,
			notes:      "notes",
			version:    1,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "file size can't be zero",
		},
		{
			name:       "zero version",
			id:         "file-id-5",
			fileName:   "test.txt",
			size:       1024,
			notes:      "notes",
			version:    0,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "version must be greater than zero",
		},
		{
			name:       "negative version",
			id:         "file-id-6",
			fileName:   "test.txt",
			size:       1024,
			notes:      "notes",
			version:    -1,
			uploadedAt: now,
			wantErr:    true,
			errMsg:     "version must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := NewFile(tt.id, tt.fileName, tt.size, tt.notes, tt.version, tt.uploadedAt)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, file)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, file)
			assert.Equal(t, tt.id, file.ID())
			assert.Equal(t, tt.fileName, file.Name())
			assert.Equal(t, tt.size, file.Size())
			assert.Equal(t, tt.notes, file.Notes())
			assert.Equal(t, tt.version, file.Version())
			assert.Equal(t, tt.uploadedAt, file.UploadedAt())
		})
	}
}

func TestFile_Getters(t *testing.T) {
	t.Parallel()

	now := time.Now()
	file, err := NewFile("test-id", "test.txt", 1024, "test notes", 1, now)
	require.NoError(t, err)

	assert.Equal(t, "test-id", file.ID())
	assert.Equal(t, "test.txt", file.Name())
	assert.Equal(t, int64(1024), file.Size())
	assert.Equal(t, "test notes", file.Notes())
	assert.Equal(t, int64(1), file.Version())
	assert.Equal(t, now, file.UploadedAt())
}
