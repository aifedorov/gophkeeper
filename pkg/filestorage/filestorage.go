package filestorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

const (
	rootPath = "storage/files/"
	dirMode  = 0700
)

type FileStorage struct {
	logger   *zap.Logger
	tmpPaths sync.Map // map[string]string: key = dirname+filename, value = tmppath
}

func NewFileStorage(logger *zap.Logger) *FileStorage {
	return &FileStorage{
		logger: logger,
	}
}

func (f *FileStorage) Upload(ctx context.Context, dirname, filename string, reader io.Reader) (path string, err error) {
	tmpFilename, err := f.createTmpFile(ctx, dirname, reader)
	if err != nil {
		return "", fmt.Errorf("filestorage: failed to create temp file: %w", err)
	}

	path = f.getFullPath(dirname, filename)
	err = f.renameTmpFile(tmpFilename, dirname, filename)
	if err != nil {
		f.logger.Error("filestorage: failed to rename temp file", zap.Error(err))
		_ = f.removeTmpFile(tmpFilename)
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
	file, err := f.openForRead(path)
	if err != nil {
		f.logger.Error("filestorage: failed to open file", zap.Error(err))
		return nil, fmt.Errorf("filestorage: failed to open file: %w", err)
	}

	return file, nil
}

func (f *FileStorage) BeginUpdate(ctx context.Context, dirname, filename string, reader io.Reader) (tmppath, targetpath string, err error) {
	f.logger.Debug("filestorage: starting file update",
		zap.String("dirname", dirname),
		zap.String("filename", filename),
	)

	tmppath, err = f.createTmpFile(ctx, dirname, reader)
	if err != nil {
		return "", "", fmt.Errorf("filestorage: failed to create temp file: %w", err)
	}

	key := filepath.Join(dirname, filename)
	f.tmpPaths.Store(key, tmppath)

	return tmppath, f.getFullPath(dirname, filename), nil
}

// CommitUpdate satisfies the binary interface signature.
// It retrieves the temp path from BeginUpdate and commits the update.
func (f *FileStorage) CommitUpdate(ctx context.Context, dirname, filename string) error {
	key := filepath.Join(dirname, filename)
	value, ok := f.tmpPaths.LoadAndDelete(key)
	if !ok {
		return fmt.Errorf("filestorage: no temp file found for %s/%s (BeginUpdate must be called first)", dirname, filename)
	}
	tmppath, ok := value.(string)
	if !ok {
		return fmt.Errorf("filestorage: invalid temp path type: %T", value)
	}

	err := f.renameTmpFile(tmppath, dirname, filename)
	if err != nil {
		f.logger.Error("filestorage: failed to rename temp file", zap.Error(err))
		return fmt.Errorf("filestorage: failed to rename temp file: %w", err)
	}

	return nil
}

func (f *FileStorage) AbortUpdate(_ context.Context, tmppath string) error {
	return f.removeTmpFile(tmppath)
}

// ReadContent reads the entire file content and returns it as a string.
// If maxSize is greater than 0 and the file exceeds this size, returns an error.
func (f *FileStorage) ReadContent(ctx context.Context, path string, maxSize int64) (string, error) {
	f.logger.Debug("filestorage: reading file content",
		zap.String("path", path),
		zap.Int64("maxSize", maxSize),
	)

	file, err := f.OpenFile(ctx, path)
	if err != nil {
		f.logger.Error("filestorage: failed to open file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

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
	return f.openForRead(path)
}
