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

type FileStorage struct {
	logger *zap.Logger
}

func NewFileStorage(logger *zap.Logger) *FileStorage {
	return &FileStorage{
		logger: logger,
	}
}

func (f *FileStorage) Upload(_ context.Context, userID, fileID string, reader io.Reader) (path string, err error) {
	f.logger.Debug("filestorage: uploading file",
		zap.String("user_id", userID),
		zap.String("file_id", fileID),
	)

	dir := f.getDir(userID)
	path = f.getFullPath(userID, fileID)

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

func (f *FileStorage) Delete(_ context.Context, userID, fileID string) error {
	f.logger.Debug("filestorage: deleting file",
		zap.String("user_id", userID),
		zap.String("file_id", fileID),
	)

	path := f.getFullPath(userID, fileID)
	return os.Remove(path)
}

func (f *FileStorage) Download(_ context.Context, userID, fileID string) (reader io.ReadCloser, err error) {
	f.logger.Debug("filestorage: getting file",
		zap.String("user_id", userID),
		zap.String("file_id", fileID),
	)

	path := f.getFullPath(userID, fileID)

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

func (f *FileStorage) Update(_ context.Context, userID, fileID string, reader io.Reader) (path string, err error) {
	f.logger.Debug("filestorage: updating file",
		zap.String("user_id", userID),
		zap.String("file_id", fileID),
	)

	path = f.getFullPath(userID, fileID)

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

func (f *FileStorage) getDir(userID string) string {
	return filepath.Join(rootPath, userID)
}

func (f *FileStorage) getFullPath(userID, fileID string) string {
	return filepath.Join(rootPath, userID, fileID)
}
