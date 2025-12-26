// Package binary provides file metadata utilities for the GophKeeper client.
package binary

import (
	"fmt"
)

// FileMeta contains metadata about a downloaded file.
// It includes the file name, size, version and optional notes.
type FileMeta struct {
	name    string
	size    int64
	notes   string
	version int64
}

// NewFileMeta creates a new FileMeta with the provided data.
// It validates that name is not empty, size and version are greater than zero.
// Returns an error if validation fails.
func NewFileMeta(name string, size int64, notes string, version int64) (*FileMeta, error) {
	if name == "" {
		return nil, fmt.Errorf("file name is required")
	}
	if size == 0 {
		return nil, fmt.Errorf("file size can't be zero")
	}
	if version < 1 {
		return nil, fmt.Errorf("version must be greater than zero")
	}

	return &FileMeta{
		name:  name,
		size:  size,
		notes: notes,
	}, nil
}

// Name returns the file's name.
func (f *FileMeta) Name() string {
	return f.name
}

// Size returns the file's size in bytes.
func (f *FileMeta) Size() int64 {
	return f.size
}

// Notes returns the file's notes/metadata.
func (f *FileMeta) Notes() string {
	return f.notes
}

// Version returns the file's version.
func (f *FileMeta) Version() int64 {
	return f.version
}
