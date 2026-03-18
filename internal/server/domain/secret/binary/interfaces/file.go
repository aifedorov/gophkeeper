// Package interfaces provides domain interfaces for binary file management.
package interfaces

import (
	"fmt"
	"time"
)

// maxFileSize is the maximum allowed file size (10GB).
const maxFileSize int64 = 10 * 1024 * 1024 * 1024

// File represents a binary file entity in the domain.
// It contains file metadata including ID, name, size, storage path, notes, and upload timestamp.
type File struct {
	id         string
	name       string
	size       int64
	path       string
	notes      string
	version    int64
	uploadedAt time.Time
}

// NewFile creates a new File entity with the provided data.
// It validates that id and name are not empty, size is greater than zero,
// and size doesn't exceed the maximum allowed size (10GB).
// Returns an error if validation fails.
func NewFile(
	id, name string,
	size int64,
	path, notes string,
	version int64,
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
	if version < 1 {
		return nil, fmt.Errorf("invalid file version: %d", version)
	}

	return &File{
		id:         id,
		name:       name,
		size:       size,
		path:       path,
		notes:      notes,
		version:    version,
		uploadedAt: uploadedAt,
	}, nil
}

// GetID returns the file's unique identifier.
func (f *File) GetID() string {
	return f.id
}

// GetName returns the file's name.
func (f *File) GetName() string {
	return f.name
}

// GetSize returns the file's size in bytes.
func (f *File) GetSize() int64 {
	return f.size
}

// GetNotes returns the file's notes/metadata.
func (f *File) GetNotes() string {
	return f.notes
}

// GetUploadedAt returns the timestamp when the file was uploaded.
func (f *File) GetUploadedAt() time.Time {
	return f.uploadedAt
}

// GetPath returns the file's storage path.
func (f *File) GetPath() string {
	return f.path
}

// SetPath sets the file's storage path.
func (f *File) SetPath(path string) {
	f.path = path
}

func (f *File) GetVersion() int64 {
	return f.version
}
