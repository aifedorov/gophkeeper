// Package interfaces provides interfaces for text note management dependencies.
package interfaces

import (
	"context"
	"io"
)

//go:generate mockgen -source=storage.go -destination=mock_storage.go -package=interfaces

// Storage defines the interface for local file storage operations needed by text service.
type Storage interface {
	// Upload creates a new file in the specified directory from a reader.
	Upload(ctx context.Context, dirname, filename string, reader io.Reader) (path string, err error)
	// ReadContent reads the entire file content and returns it as a string.
	// If maxSize is greater than 0 and the file exceeds this size, returns an error.
	ReadContent(ctx context.Context, path string, maxSize int64) (string, error)
}
