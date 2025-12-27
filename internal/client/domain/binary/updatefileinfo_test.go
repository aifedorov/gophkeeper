package binary

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUpdateFileInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		notes     string
		wantErr   bool
		wantNotes string
	}{
		{
			name:      "successful creation",
			id:        "file-123",
			notes:     "test notes",
			wantErr:   false,
			wantNotes: "test notes",
		},
		{
			name:      "empty notes",
			id:        "file-456",
			notes:     "",
			wantErr:   false,
			wantNotes: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a temp file for testing
			tmpFile, err := os.CreateTemp("", "updatefileinfo_test_*.txt")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			// Write some content to get a non-zero size
			_, err = tmpFile.WriteString("test content")
			require.NoError(t, err)

			info, err := NewUpdateFileInfo(tt.id, tmpFile, tt.notes)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, info)
			} else {
				require.NoError(t, err)
				require.NotNil(t, info)
				assert.Equal(t, tt.id, info.ID())
				assert.NotEmpty(t, info.Name())
				assert.Equal(t, int64(12), info.Size()) // "test content" is 12 bytes
				assert.Equal(t, tt.wantNotes, info.Notes())
				assert.Equal(t, int64(0), info.Version()) // default version is 0
			}
		})
	}
}

func TestNewUpdateFileInfoWithVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          string
		notes       string
		version     int64
		wantErr     bool
		wantNotes   string
		wantVersion int64
	}{
		{
			name:        "successful creation with version",
			id:          "file-123",
			notes:       "test notes",
			version:     5,
			wantErr:     false,
			wantNotes:   "test notes",
			wantVersion: 5,
		},
		{
			name:        "empty notes with version",
			id:          "file-456",
			notes:       "",
			version:     10,
			wantErr:     false,
			wantNotes:   "",
			wantVersion: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a temp file for testing
			tmpFile, err := os.CreateTemp("", "updatefileinfo_version_test_*.txt")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			// Write some content to get a non-zero size
			_, err = tmpFile.WriteString("test content")
			require.NoError(t, err)

			info, err := NewUpdateFileInfoWithVersion(tt.id, tmpFile, tt.notes, tt.version)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, info)
			} else {
				require.NoError(t, err)
				require.NotNil(t, info)
				assert.Equal(t, tt.id, info.ID())
				assert.NotEmpty(t, info.Name())
				assert.Equal(t, int64(12), info.Size()) // "test content" is 12 bytes
				assert.Equal(t, tt.wantNotes, info.Notes())
				assert.Equal(t, tt.wantVersion, info.Version())
			}
		})
	}
}

func TestNewUpdateFileInfo_StatError(t *testing.T) {
	t.Parallel()

	// Create a temp file, close it, then try to stat (should fail on closed file)
	tmpFile, err := os.CreateTemp("", "updatefileinfo_stat_error_*.txt")
	require.NoError(t, err)
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	// Close the file first
	tmpFile.Close()

	// Try to create UpdateFileInfo with closed file - Stat() will fail
	info, err := NewUpdateFileInfo("id", tmpFile, "notes")
	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Contains(t, err.Error(), "failed to get file stat")
}

func TestNewUpdateFileInfoWithVersion_StatError(t *testing.T) {
	t.Parallel()

	// Create a temp file, close it, then try to stat (should fail on closed file)
	tmpFile, err := os.CreateTemp("", "updatefileinfo_version_stat_error_*.txt")
	require.NoError(t, err)
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	// Close the file first
	tmpFile.Close()

	// Try to create UpdateFileInfo with closed file - Stat() will fail
	info, err := NewUpdateFileInfoWithVersion("id", tmpFile, "notes", 1)
	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Contains(t, err.Error(), "failed to get file stat")
}

func TestUpdateFileInfo_Getters(t *testing.T) {
	t.Parallel()

	// Create a temp file for testing
	tmpFile, err := os.CreateTemp("", "updatefileinfo_getters_*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString("hello world")
	require.NoError(t, err)

	info, err := NewUpdateFileInfoWithVersion("test-id", tmpFile, "my notes", 7)
	require.NoError(t, err)

	// Test ID getter
	assert.Equal(t, "test-id", info.ID())

	// Test Name getter
	assert.NotEmpty(t, info.Name())

	// Test Size getter
	assert.Equal(t, int64(11), info.Size())

	// Test Notes getter
	assert.Equal(t, "my notes", info.Notes())

	// Test Version getter
	assert.Equal(t, int64(7), info.Version())
}
