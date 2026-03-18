// Package binary provides client interfaces for binary file management.
package binary

import (
	"context"
	"io"
)

//go:generate mockgen -source=interface.go -destination=mock_client_test.go -package=binary

// Client defines the interface for gRPC client operations for binary file management.
// This interface abstracts the gRPC communication layer.
type Client interface {
	// Upload sends a file to the server using streaming.
	// The fileInfo contains metadata, and reader provides the file content.
	// Returns the file ID and version number after successful upload.
	Upload(ctx context.Context, fileInfo *FileInfo, reader io.Reader) (id string, version int64, err error)
	// List retrieves all files for the authenticated user from the server.
	List(ctx context.Context) ([]File, error)
	// Download retrieves a file by ID from the server using streaming.
	// Returns a ReadCloser for the file content and metadata. The reader should be closed after use.
	Download(ctx context.Context, id string) (io.ReadCloser, *FileMeta, error)
	// Update sends an updated file to the server using streaming.
	// The fileInfo contains metadata including file ID, and reader provides the new file content.
	// Returns the new version number after successful update.
	Update(ctx context.Context, fileInfo *UpdateFileInfo, reader io.Reader) (version int64, err error)
	// Delete sends a request to delete a file by ID from the server.
	Delete(ctx context.Context, id string) error
}
