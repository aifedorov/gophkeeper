package filestorage

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewFileStorage(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	fs := NewFileStorage(logger)

	require.NotNil(t, fs)
	assert.Equal(t, logger, fs.logger)
}

func TestFileStorage_Upload(t *testing.T) {
	t.Parallel()

	// Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", "filestorage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	fs := NewFileStorage(logger)

	tests := []struct {
		name     string
		dirname  string
		filename string
		content  string
		wantErr  bool
	}{
		{
			name:     "successful upload",
			dirname:  filepath.Join(tempDir, "user1"),
			filename: "test.txt",
			content:  "hello world",
			wantErr:  false,
		},
		{
			name:     "upload to new directory",
			dirname:  filepath.Join(tempDir, "user2/subdir"),
			filename: "file.txt",
			content:  "test content",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			reader := strings.NewReader(tt.content)

			path, err := fs.Upload(ctx, tt.dirname, tt.filename, reader)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, path)

				// Verify file exists and content is correct
				content, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.Equal(t, tt.content, string(content))

				// Cleanup
				os.Remove(path)
			}
		})
	}
}

func TestFileStorage_Delete(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	fs := NewFileStorage(logger)

	t.Run("successful delete", func(t *testing.T) {
		ctx := context.Background()
		dirname := "delete_test_user"
		filename := "delete_test.txt"

		// First upload a file
		reader := strings.NewReader("content to delete")
		path, err := fs.Upload(ctx, dirname, filename, reader)
		require.NoError(t, err)

		// Now delete it
		err = fs.Delete(ctx, dirname, filename)
		require.NoError(t, err)

		// Verify file is deleted
		_, err = os.Stat(path)
		assert.True(t, os.IsNotExist(err))

		// Cleanup directory
		os.RemoveAll(filepath.Join("storage/files", dirname))
	})

	t.Run("delete non-existent file", func(t *testing.T) {
		ctx := context.Background()

		err := fs.Delete(ctx, "nonexistent_dir", "nonexistent.txt")
		assert.Error(t, err)
	})
}

func TestFileStorage_Download(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	fs := NewFileStorage(logger)

	t.Run("successful download", func(t *testing.T) {
		ctx := context.Background()
		dirname := "download_test_user"
		filename := "download_test.txt"
		testContent := "test file content for download"

		// First upload a file
		reader := strings.NewReader(testContent)
		_, err := fs.Upload(ctx, dirname, filename, reader)
		require.NoError(t, err)
		defer os.RemoveAll(filepath.Join("storage/files", dirname))

		// Now download it
		downloadReader, err := fs.Download(ctx, dirname, filename)
		require.NoError(t, err)
		defer downloadReader.Close()

		// Read content
		buf := make([]byte, len(testContent))
		_, err = downloadReader.Read(buf)
		require.NoError(t, err)
		assert.Equal(t, testContent, string(buf))
	})

	t.Run("download non-existent file", func(t *testing.T) {
		ctx := context.Background()

		reader, err := fs.Download(ctx, "nonexistent_dir", "nonexistent.txt")
		assert.Error(t, err)
		assert.Nil(t, reader)
	})
}

func TestFileStorage_BeginUpdate_CommitUpdate_AbortUpdate(t *testing.T) {
	t.Parallel()

	// Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", "filestorage_update_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	fs := NewFileStorage(logger)

	t.Run("begin and commit update", func(t *testing.T) {
		ctx := context.Background()
		dirname := filepath.Join(tempDir, "user1")
		filename := "update_test.txt"
		content := "updated content"

		reader := strings.NewReader(content)
		tmppath, targetpath, err := fs.BeginUpdate(ctx, dirname, filename, reader)
		require.NoError(t, err)
		assert.NotEmpty(t, tmppath)
		assert.NotEmpty(t, targetpath)

		// Commit the update
		err = fs.CommitUpdate(ctx, dirname, filename)
		require.NoError(t, err)

		// Verify file exists at target path
		fileContent, err := os.ReadFile(targetpath)
		require.NoError(t, err)
		assert.Equal(t, content, string(fileContent))

		// Cleanup
		os.Remove(targetpath)
	})

	t.Run("begin and abort update", func(t *testing.T) {
		ctx := context.Background()
		dirname := filepath.Join(tempDir, "user2")
		filename := "abort_test.txt"
		content := "content to abort"

		reader := strings.NewReader(content)
		tmppath, _, err := fs.BeginUpdate(ctx, dirname, filename, reader)
		require.NoError(t, err)

		// Abort the update
		err = fs.AbortUpdate(ctx, tmppath)
		require.NoError(t, err)

		// Verify temp file is deleted
		_, err = os.Stat(tmppath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("commit without begin", func(t *testing.T) {
		ctx := context.Background()

		err := fs.CommitUpdate(ctx, tempDir, "no_begin.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no temp file found")
	})
}

func TestFileStorage_ReadContent(t *testing.T) {
	t.Parallel()

	// Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", "filestorage_read_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	fs := NewFileStorage(logger)

	t.Run("successful read", func(t *testing.T) {
		// Create a test file
		testContent := "test file content for reading"
		testFile := filepath.Join(tempDir, "read_test.txt")
		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))

		ctx := context.Background()

		content, err := fs.ReadContent(ctx, testFile, 0)
		require.NoError(t, err)
		assert.Equal(t, testContent, content)
	})

	t.Run("read with max size within limit", func(t *testing.T) {
		testContent := "small content"
		testFile := filepath.Join(tempDir, "small_read_test.txt")
		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))

		ctx := context.Background()

		content, err := fs.ReadContent(ctx, testFile, 100)
		require.NoError(t, err)
		assert.Equal(t, testContent, content)
	})

	t.Run("read with max size exceeded", func(t *testing.T) {
		testContent := "this is a longer content that should exceed the limit"
		testFile := filepath.Join(tempDir, "large_read_test.txt")
		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))

		ctx := context.Background()

		content, err := fs.ReadContent(ctx, testFile, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds maximum")
		assert.Empty(t, content)
	})

	t.Run("read non-existent file", func(t *testing.T) {
		ctx := context.Background()

		content, err := fs.ReadContent(ctx, "/nonexistent/file.txt", 0)
		assert.Error(t, err)
		assert.Empty(t, content)
	})
}

func TestFileStorage_OpenFile(t *testing.T) {
	t.Parallel()

	// Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", "filestorage_openfile_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	fs := NewFileStorage(logger)

	t.Run("successful open", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "openfile_test.txt")
		require.NoError(t, os.WriteFile(testFile, []byte("test"), 0644))

		ctx := context.Background()

		file, err := fs.OpenFile(ctx, testFile)
		require.NoError(t, err)
		defer file.Close()

		assert.NotNil(t, file)
	})

	t.Run("open non-existent file", func(t *testing.T) {
		ctx := context.Background()

		file, err := fs.OpenFile(ctx, "/nonexistent/file.txt")
		assert.Error(t, err)
		assert.Nil(t, file)
	})
}
