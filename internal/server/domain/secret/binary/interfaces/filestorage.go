// Package interfaces provides file storage interfaces for binary file management.
package interfaces

import (
	"context"
	"io"
)

//go:generate mockgen -source=filestorage.go -destination=mocks/mock_filestorage.go -package=mocks

// FileStorage defines the interface for physical file storage operations.
// Files are stored encrypted and organized by user ID and file ID.
type FileStorage interface {
	// Upload stores a file from the provided reader for the specified user and file ID.
	// Returns the file path where the file was stored.
	Upload(ctx context.Context, userID, fileID string, reader io.Reader) (filepath string, err error)
	// Delete removes a physical file for the specified user and file ID.
	Delete(ctx context.Context, userID, fileID string) error
	// Download retrieves a physical file for the specified user and file ID.
	// Returns a ReadCloser that should be closed by the caller after use.
	Download(_ context.Context, userID, fileID string) (reader io.ReadCloser, err error)
}
