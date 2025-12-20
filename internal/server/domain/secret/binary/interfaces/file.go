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
	path       string
	notes      string
	uploadedAt time.Time
}

func NewFile(
	id, name string,
	size int64,
	path, notes string,
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

	return &File{
		id:         id,
		name:       name,
		size:       size,
		path:       path,
		notes:      notes,
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

func (f *File) GetNotes() string {
	return f.notes
}

func (f *File) GetUploadedAt() time.Time {
	return f.uploadedAt
}

func (f *File) GetPath() string {
	return f.path
}

func (f *File) SetPath(path string) {
	f.path = path
}
