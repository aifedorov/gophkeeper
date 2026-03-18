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

func (f *FileStorage) getDir(dirname string) string {
	return filepath.Join(f.rootPath, dirname)
}

func (f *FileStorage) getFullPath(dirname, filename string) string {
	return filepath.Join(f.rootPath, dirname, filename)
}

func (f *FileStorage) createTmpFile(_ context.Context, dirname string, reader io.Reader) (filename string, err error) {
	f.logger.Debug("filestorage: creating temp file for upload",
		zap.String("dirname", dirname),
	)

	dir := f.getDir(dirname)
	if err := os.MkdirAll(dir, dirMode); err != nil {
		f.logger.Error("filestorage: failed to create directory", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to create directory: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, "upload-*.tmp")
	if err != nil {
		f.logger.Error("filestorage: failed to create temp file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to create temp file: %w", err)
	}

	tmpPath := tmpFile.Name()
	defer func() {
		if tmpFile != nil {
			_ = tmpFile.Close()
		}
		if err != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	f.logger.Debug("filestorage: created temp file", zap.String("path", tmpPath))

	_, err = io.Copy(tmpFile, reader)
	if err != nil {
		f.logger.Error("filestorage: failed to upload file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to upload file: %w", err)
	}

	f.logger.Debug("filestorage: file copied successfully", zap.String("path", tmpPath))

	err = tmpFile.Close()
	if err != nil {
		f.logger.Error("filestorage: failed to close temp file", zap.Error(err))
		return "", fmt.Errorf("filestorage: failed to close temp file: %w", err)
	}

	return tmpPath, nil
}

func (f *FileStorage) renameTmpFile(tmppath, dirname, filename string) error {
	err := os.Rename(tmppath, f.getFullPath(dirname, filename))
	if err != nil {
		f.logger.Error("filestorage: failed to rename temp file", zap.Error(err))
		return fmt.Errorf("filestorage: failed to rename temp file: %w", err)
	}
	return nil
}

func (f *FileStorage) removeTmpFile(tmppath string) error {
	f.logger.Debug("filestorage: removing temp file", zap.String("path", tmppath))
	return os.Remove(tmppath)
}

func (f *FileStorage) openForRead(path string) (*os.File, error) {
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
