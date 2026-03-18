package binary

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		notes     string
		wantErr   bool
		wantNotes string
	}{
		{
			name:      "successful creation",
			notes:     "test notes",
			wantErr:   false,
			wantNotes: "test notes",
		},
		{
			name:      "empty notes",
			notes:     "",
			wantErr:   false,
			wantNotes: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a temp file for testing
			tmpFile, err := os.CreateTemp("", "fileinfo_test_*.txt")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			// Write some content to get a non-zero size
			_, err = tmpFile.WriteString("test content")
			require.NoError(t, err)

			info, err := NewFileInfo(tmpFile, tt.notes)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, info)
			} else {
				require.NoError(t, err)
				require.NotNil(t, info)
				assert.NotEmpty(t, info.Name())
				assert.Equal(t, int64(12), info.Size()) // "test content" is 12 bytes
				assert.Equal(t, tt.wantNotes, info.Notes())
			}
		})
	}
}

func TestNewFileInfo_StatError(t *testing.T) {
	t.Parallel()

	// Create a temp file, close it, then try to stat (should fail on closed file)
	tmpFile, err := os.CreateTemp("", "fileinfo_stat_error_*.txt")
	require.NoError(t, err)
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	// Close the file first
	tmpFile.Close()

	// Try to create FileInfo with closed file - Stat() will fail
	info, err := NewFileInfo(tmpFile, "notes")
	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Contains(t, err.Error(), "failed to get file stat")
}

func TestFileInfo_Getters(t *testing.T) {
	t.Parallel()

	// Create a temp file for testing
	tmpFile, err := os.CreateTemp("", "fileinfo_getters_*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString("hello world")
	require.NoError(t, err)

	info, err := NewFileInfo(tmpFile, "my notes")
	require.NoError(t, err)

	// Test Name getter
	assert.NotEmpty(t, info.Name())

	// Test Size getter
	assert.Equal(t, int64(11), info.Size())

	// Test Notes getter
	assert.Equal(t, "my notes", info.Notes())
}
