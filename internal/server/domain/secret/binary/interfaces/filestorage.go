// Package interfaces provides file storage interfaces for binary file management.
package interfaces

import (
	"context"
	"io"
)

//go:generate mockgen -source=filestorage.go -destination=mocks/mock_filestorage.go -package=mocks

// FileStorage defines the interface for physical file storage operations.
// Files are stored encrypted and organized by directory and filename.
type FileStorage interface {
	// Upload stores a file from the provided reader for the specified directory and filename.
	Upload(ctx context.Context, dirname, filename string, reader io.Reader) (filepath string, err error)
	// Delete removes a physical file for the specified directory and filename.
	Delete(ctx context.Context, dirname, filename string) error
	// Download retrieves a physical file for the specified directory and filename.
	// Returns a ReadCloser that should be closed by the caller after use.
	Download(ctx context.Context, dirname, filename string) (reader io.ReadCloser, err error)

	BeginUpdate(ctx context.Context, dirname, filename string, reader io.Reader) (tmppath, targetpath string, err error)

	CommitUpdate(ctx context.Context, dirname, filename string) error

	AbortUpdate(ctx context.Context, tmppath string) error
}
