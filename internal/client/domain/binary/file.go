package binary

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
)

type FileInfo struct {
	name     string
	size     int64
	mimeType string
}

func NewFileInfo(file *os.File) (*FileInfo, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	mimeType, err := extractMimeType(file.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to extract mime type: %w", err)
	}

	return &FileInfo{
		name:     stat.Name(),
		size:     stat.Size(),
		mimeType: mimeType,
	}, nil
}

func extractMimeType(filename string) (string, error) {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "", fmt.Errorf("failed to extract mime type: %s", ext)
	}
	return mimeType, nil
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Size() int64 {
	return f.size
}

func (f *FileInfo) MimeType() string {
	return f.mimeType
}
