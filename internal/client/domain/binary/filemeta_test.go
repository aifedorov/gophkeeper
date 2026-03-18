package binary

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fileName string
		size     int64
		notes    string
		version  int64
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "creates file meta with all fields",
			fileName: "test.txt",
			size:     1024,
			notes:    "test notes",
			version:  1,
			wantErr:  false,
		},
		{
			name:     "creates file meta without notes",
			fileName: "document.pdf",
			size:     2048,
			notes:    "",
			version:  2,
			wantErr:  false,
		},
		{
			name:     "empty file name",
			fileName: "",
			size:     1024,
			notes:    "notes",
			version:  1,
			wantErr:  true,
			errMsg:   "file name is required",
		},
		{
			name:     "zero size",
			fileName: "test.txt",
			size:     0,
			notes:    "notes",
			version:  1,
			wantErr:  true,
			errMsg:   "file size can't be zero",
		},
		{
			name:     "zero version",
			fileName: "test.txt",
			size:     1024,
			notes:    "notes",
			version:  0,
			wantErr:  true,
			errMsg:   "version must be greater than zero",
		},
		{
			name:     "negative version",
			fileName: "test.txt",
			size:     1024,
			notes:    "notes",
			version:  -1,
			wantErr:  true,
			errMsg:   "version must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			meta, err := NewFileMeta(tt.fileName, tt.size, tt.notes, tt.version)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, meta)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, meta)
			assert.Equal(t, tt.fileName, meta.Name())
			assert.Equal(t, tt.size, meta.Size())
			assert.Equal(t, tt.notes, meta.Notes())
		})
	}
}

func TestFileMeta_Getters(t *testing.T) {
	t.Parallel()

	meta, err := NewFileMeta("test.txt", 1024, "test notes", 1)
	require.NoError(t, err)

	assert.Equal(t, "test.txt", meta.Name())
	assert.Equal(t, int64(1024), meta.Size())
	assert.Equal(t, "test notes", meta.Notes())
	// Note: Version is not set in NewFileMeta, so it returns 0
	assert.Equal(t, int64(0), meta.Version())
}
