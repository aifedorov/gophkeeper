// Package interfaces provides interfaces for binary file management dependencies.
package interfaces

import (
	"context"
	"io"
	"os"
)

//go:generate mockgen -source=storage.go -destination=mock_storage.go -package=interfaces

// Storage defines the interface for local file storage operations needed by binary service.
type Storage interface {
	// Upload creates a new file in the specified directory from a reader.
	Upload(ctx context.Context, dirname, filename string, reader io.Reader) (path string, err error)
	// OpenFile opens a file for reading and returns the file handle.
	OpenFile(ctx context.Context, path string) (*os.File, error)
}
