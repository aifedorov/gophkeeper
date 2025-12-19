package interfaces

import (
	"fmt"
	"time"
)

const maxFileSize int64 = 10 * 1024 * 1024 * 1024

type File struct {
	id         string
	name       string
	size       int64
	mimeType   string
	uploadedAt time.Time
}

func NewFile(
	id, name string,
	size int64,
	mimeType string,
	uploadedAt time.Time,
) (*File, error) {
	if id == "" {
		return nil, fmt.Errorf("file id is required")
	}
	if name == "" {
		return nil, fmt.Errorf("file name is required")
	}
	if size == 0 {
		return nil, fmt.Errorf("file size is required")
	}
	if size > maxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size: %d", maxFileSize)
	}
	if mimeType == "" {
		return nil, fmt.Errorf("file mime type is required")
	}

	return &File{
		id:         id,
		name:       name,
		size:       size,
		mimeType:   mimeType,
		uploadedAt: uploadedAt,
	}, nil
}

func (f *File) GetID() string {
	return f.id
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) GetSize() int64 {
	return f.size
}

func (f *File) GetMimeType() string {
	return f.mimeType
}

func (f *File) GetUploadedAt() time.Time {
	return f.uploadedAt
}
