// Package binary provides binary file domain entities for the GophKeeper client.
package binary

import (
	"time"
)

// File represents a binary file entity in the client domain.
// It contains file metadata including ID, name, size, notes, and upload timestamp.
type File struct {
	id         string
	name       string
	size       int64
	notes      string
	uploadedAt time.Time
}

// NewFile creates a new File entity with the provided data.
// Returns an error if validation fails (currently always returns nil).
func NewFile(id, name string, size int64, notes string, uploadedAt time.Time) (*File, error) {
	return &File{id: id, name: name, size: size, notes: notes, uploadedAt: uploadedAt}, nil
}

// ID returns the file's unique identifier.
func (f *File) ID() string {
	return f.id
}

// Name returns the file's name.
func (f *File) Name() string {
	return f.name
}

// Size returns the file's size in bytes.
func (f *File) Size() int64 {
	return f.size
}

// Notes returns the file's notes/metadata.
func (f *File) Notes() string {
	return f.notes
}

// UploadedAt returns the timestamp when the file was uploaded.
func (f *File) UploadedAt() time.Time {
	return f.uploadedAt
}
