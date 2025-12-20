package binary

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
)

type FileMeta struct {
	name     string
	size     int64
	mimeType string
	notes    string
}

func NewFileMeta(file *os.File, notes string) (*FileMeta, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	mimeType, err := extractMimeType(file.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to extract mime type: %w", err)
	}

	return &FileMeta{
		name:     stat.Name(),
		size:     stat.Size(),
		mimeType: mimeType,
		notes:    notes,
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

func (f *FileMeta) Name() string {
	return f.name
}

func (f *FileMeta) Size() int64 {
	return f.size
}

func (f *FileMeta) MimeType() string {
	return f.mimeType
}

func (f *FileMeta) Notes() string {
	return f.notes
}
