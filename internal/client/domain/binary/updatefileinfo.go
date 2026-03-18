// Package binary provides file information utilities for the GophKeeper client.
package binary

import (
	"fmt"
	"os"
)

// UpdateFileInfo contains metadata about a file to be updated.
// It includes the file ID, name, size, version and optional notes.
type UpdateFileInfo struct {
	id      string
	name    string
	size    int64
	notes   string
	version int64
}

// NewUpdateFileInfo creates a new UpdateFileInfo from an open file, ID, and notes.
// It reads the file's stat information to get the name and size.
// Returns an error if the file stat cannot be read.
func NewUpdateFileInfo(id string, file *os.File, notes string) (*UpdateFileInfo, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	return &UpdateFileInfo{
		id:    id,
		name:  stat.Name(),
		size:  stat.Size(),
		notes: notes,
	}, nil
}

// NewUpdateFileInfoWithVersion creates a new UpdateFileInfo from an open file, ID, notes, and version.
// It reads the file's stat information to get the name and size.
// Returns an error if the file stat cannot be read.
func NewUpdateFileInfoWithVersion(id string, file *os.File, notes string, version int64) (*UpdateFileInfo, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	return &UpdateFileInfo{
		id:      id,
		name:    stat.Name(),
		size:    stat.Size(),
		notes:   notes,
		version: version,
	}, nil
}

// ID returns the file's ID.
func (f *UpdateFileInfo) ID() string {
	return f.id
}

// Name returns the file's name.
func (f *UpdateFileInfo) Name() string {
	return f.name
}

// Size returns the file's size in bytes.
func (f *UpdateFileInfo) Size() int64 {
	return f.size
}

// Notes returns the file's notes/metadata.
func (f *UpdateFileInfo) Notes() string {
	return f.notes
}

// Version returns the file's version.
func (f *UpdateFileInfo) Version() int64 {
	return f.version
}
