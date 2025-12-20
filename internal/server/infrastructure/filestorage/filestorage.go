package filestorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/binary/interfaces"
	"go.uber.org/zap"
)

const (
	rootPath = "storage/files/"
	dirMode  = 0700
)

type fileStorage struct {
	logger *zap.Logger
}

func NewFileStorage(logger *zap.Logger) interfaces.FileStorage {
	return &fileStorage{
		logger: logger,
	}
}

func (f *fileStorage) Upload(_ context.Context, userID, fileID string, reader io.Reader) (path string, err error) {
	f.logger.Debug("filestorage: uploading file",
		zap.String("user_id", userID),
		zap.String("file_id", fileID),
	)

	dir := filepath.Join(rootPath, userID)
	path = filepath.Join(dir, fileID)

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

func (f *fileStorage) Delete(_ context.Context, userID, fileID string) error {
	f.logger.Debug("filestorage: deleting file",
		zap.String("user_id", userID),
		zap.String("file_id", fileID),
	)
	return os.Remove(filepath.Join(rootPath, userID, fileID))
}
