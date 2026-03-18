// Package binary provides file information utilities for the GophKeeper client.
package binary

import (
	"fmt"
	"os"
)

// FileInfo contains metadata about a file to be uploaded.
// It includes the file name, size, and optional notes.
type FileInfo struct {
	name  string
	size  int64
	notes string
}

// NewFileInfo creates a new FileInfo from an open file and notes.
// It reads the file's stat information to get the name and size.
// Returns an error if the file stat cannot be read.
func NewFileInfo(file *os.File, notes string) (*FileInfo, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	return &FileInfo{
		name:  stat.Name(),
		size:  stat.Size(),
		notes: notes,
	}, nil
}

// Name returns the file's name.
func (f *FileInfo) Name() string {
	return f.name
}

// Size returns the file's size in bytes.
func (f *FileInfo) Size() int64 {
	return f.size
}

// Notes returns the file's notes/metadata.
func (f *FileInfo) Notes() string {
	return f.notes
}
