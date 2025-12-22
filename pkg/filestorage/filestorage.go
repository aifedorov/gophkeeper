package filestorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const (
	rootPath = "storage/files/"
	dirMode  = 0700
)

//go:generate mockgen -destination=mock_storage.go -package=filestorage github.com/aifedorov/gophkeeper/pkg/filestorage Storage

// Storage defines the interface for file storage operations.
type Storage interface {
	// Upload creates a new file in the specified directory from a reader.
	Upload(ctx context.Context, dirname, filename string, reader io.Reader) (path string, err error)
	// Delete removes a file from the specified directory.
	Delete(ctx context.Context, dirname, filename string) error
	// Download opens a file for reading from the specified directory.
	Download(ctx context.Context, dirname, filename string) (reader io.ReadCloser, err error)
	// Update replaces the content of an existing file with data from a reader.
	Update(ctx context.Context, dirname, filename string, reader io.Reader) (path string, err error)
	// ReadContent reads the entire file content and returns it as a string.
	// If maxSize is greater than 0 and the file exceeds this size, returns an error.
	ReadContent(ctx context.Context, path string, maxSize int64) (string, error)
	// OpenFile opens a file for reading and returns the file handle.
	// The caller is responsible for closing the file.
	OpenFile(ctx context.Context, path string) (*os.File, error)
}

type FileStorage struct {
	logger *zap.Logger
}

func NewFileStorage(logger *zap.Logger) *FileStorage {
	return &FileStorage{
		logger: logger,
	}
}

func (f *FileStorage) Upload(_ context.Context, dirname, filename string, reader io.Reader) (path string, err error) {
	f.logger.Debug("filestorage: uploading file",
		zap.String("dirname", dirname),
		zap.String("filename", filename),
	)

	dir := f.getDir(dirname)
	path = filepath.Join(dir, filename)

	if err := os.MkdirAll(dir, dirMode); err != nil {
		f.logger.Error("filestorage: failed to create directory", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to create directory: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, "upload-*.tmp")
	defer func() {
		if err != nil && tmpFile != nil {
			_ = tmpFile.Close()
			_ = os.Remove(tmpFile.Name())
		}
	}()
	if err != nil {
		f.logger.Error("filestorage: failed to create temp file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to create temp file: %w", err)
	}

	f.logger.Debug("filestorage: created temp file", zap.String("path", tmpFile.Name()))

	_, err = io.Copy(tmpFile, reader)
	if err != nil {
		f.logger.Error("filestorage: failed to upload file", zap.Error(err))
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("filestorage: failed to upload file: %w", err)
	}

	f.logger.Debug("filestorage: file copied successfully", zap.String("path", tmpFile.Name()))

	err = tmpFile.Close()
	if err != nil {
		f.logger.Error("filestorage: failed to close temp file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to close temp file: %w", err)
	}

	err = os.Rename(tmpFile.Name(), path)
	if err != nil {
		f.logger.Error("filestorage: failed to rename temp file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to rename temp file: %w", err)
	}

	f.logger.Debug("filestorage: file uploaded successfully", zap.String("path", path))
	return path, nil
}

func (f *FileStorage) Delete(_ context.Context, dirname, filename string) error {
	f.logger.Debug("filestorage: deleting file",
		zap.String("dirname", dirname),
		zap.String("filename", filename),
	)

	path := f.getFullPath(dirname, filename)
	return os.Remove(path)
}

func (f *FileStorage) Download(_ context.Context, dirname, filename string) (reader io.ReadCloser, err error) {
	f.logger.Debug("filestorage: getting file",
		zap.String("dirname", dirname),
		zap.String("filename", filename),
	)

	path := f.getFullPath(dirname, filename)

	// #nosec G304
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		f.logger.Debug("filestorage: file not found", zap.String("path", path))
		return nil, fmt.Errorf("filestorage: file not found: %w", err)
	}
	if err != nil {
		f.logger.Error("filestorage: failed to open file", zap.Error(err))
		return nil, fmt.Errorf("filestorage: failed to open file: %w", err)
	}
	return file, nil
}

func (f *FileStorage) Update(_ context.Context, dirname, filename string, reader io.Reader) (path string, err error) {
	f.logger.Debug("filestorage: updating file",
		zap.String("dirname", dirname),
		zap.String("filename", filename),
	)

	path = f.getFullPath(dirname, filename)

	// #nosec G304
	file, err := os.Create(path)
	defer func() {
		if err != nil && file != nil {
			_ = file.Close()
		}
	}()
	if err != nil {
		f.logger.Error("filestorage: failed to create file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to create file: %w", err)
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		f.logger.Error("filestorage: failed to update file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to update file: %w", err)
	}

	err = file.Close()
	if err != nil {
		f.logger.Error("filestorage: failed to close file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to close file: %w", err)
	}

	f.logger.Debug("filestorage: file updated successfully", zap.String("path", path))
	return path, nil
}

// ReadContent reads the entire file content and returns it as a string.
// If maxSize is greater than 0 and the file exceeds this size, returns an error.
func (f *FileStorage) ReadContent(ctx context.Context, path string, maxSize int64) (string, error) {
	f.logger.Debug("filestorage: reading file content",
		zap.String("path", path),
		zap.Int64("maxSize", maxSize),
	)

	// #nosec G304
	file, err := os.Open(path)
	if err != nil {
		f.logger.Error("filestorage: failed to open file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if maxSize > 0 {
		fileInfo, err := file.Stat()
		if err != nil {
			f.logger.Error("filestorage: failed to get file info", zap.Error(err))
			return "", fmt.Errorf("filestorage: failed to get file info: %w", err)
		}
		if fileInfo.Size() > maxSize {
			f.logger.Debug("filestorage: file exceeds max size",
				zap.Int64("fileSize", fileInfo.Size()),
				zap.Int64("maxSize", maxSize),
			)
			return "", fmt.Errorf("filestorage: file size %d exceeds maximum %d", fileInfo.Size(), maxSize)
		}
	}

	content, err := io.ReadAll(file)
	if err != nil {
		f.logger.Error("filestorage: failed to read file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to read file: %w", err)
	}

	f.logger.Debug("filestorage: file content read successfully",
		zap.Int("contentSize", len(content)),
	)
	return string(content), nil
}

// OpenFile opens a file for reading and returns the file handle.
// The caller is responsible for closing the file.
func (f *FileStorage) OpenFile(ctx context.Context, path string) (*os.File, error) {
	f.logger.Debug("filestorage: opening file",
		zap.String("path", path),
	)

	// #nosec G304
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		f.logger.Debug("filestorage: file not found", zap.String("path", path))
		return nil, fmt.Errorf("filestorage: file not found: %w", err)
	}
	if err != nil {
		f.logger.Error("filestorage: failed to open file", zap.Error(err))
		return nil, fmt.Errorf("filestorage: failed to open file: %w", err)
	}

	f.logger.Debug("filestorage: file opened successfully", zap.String("path", path))
	return file, nil
}

func (f *FileStorage) getDir(dirname string) string {
	return filepath.Join(rootPath, dirname)
}

func (f *FileStorage) getFullPath(dirname, filename string) string {
	return filepath.Join(rootPath, dirname, filename)
}
