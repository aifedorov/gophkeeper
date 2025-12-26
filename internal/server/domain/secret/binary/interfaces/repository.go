// Package interfaces provides repository interfaces for binary file storage.
package interfaces

import (
	"context"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock_repository.go -package=mocks

// Repository defines the interface for binary file repository operations.
// All operations are scoped to a specific user ID.
type Repository interface {
	// Create stores a new file record in the repository for the specified user.
	Create(ctx context.Context, userID string, file RepositoryFile) (*RepositoryFile, error)
	// Get retrieves a file record by ID for the specified user.
	// Returns an error if the file doesn't exist or doesn't belong to the user.
	Get(ctx context.Context, userID, id string) (*RepositoryFile, error)
	// List retrieves all file records for the specified user.
	List(ctx context.Context, userID string) ([]RepositoryFile, error)
	// Update updates a file record by ID for the specified user.
	// Returns an error if the file doesn't exist or doesn't belong to the user.
	Update(ctx context.Context, userID, id string, file RepositoryFile) (*RepositoryFile, error)
	// Delete removes a file record by ID for the specified user.
	// Returns an error if the file doesn't exist or doesn't belong to the user.
	Delete(ctx context.Context, userID, id string) error
}
