// Package credential provides client interfaces for credential management.
package credential

import (
	"context"
)

// Client defines the interface for gRPC client operations for credential management.
// This interface abstracts the gRPC communication layer.
type Client interface {
	// Create sends a request to create a new credential on the server.
	Create(ctx context.Context, creds Credential) (id string, version int64, err error)
	// Update sends a request to update an existing credential on the server.
	// Returns the new version number after successful update.
	Update(ctx context.Context, id string, creds Credential) (version int64, err error)
	// Delete sends a request to delete a credential by ID from the server.
	Delete(ctx context.Context, id string) error
	// List retrieves all credentials for the authenticated user from the server.
	List(ctx context.Context) ([]Credential, error)
}

// CacheStorage defines the interface for caching credential version numbers locally.
// This is used for optimistic locking to detect concurrent modifications.
type CacheStorage interface {
	// SetCredentialVersion stores the version number for a credential in the cache.
	SetCredentialVersion(id string, version int64) error
	// GetCredentialVersion retrieves the cached version number for a credential.
	// Returns an error if the credential is not found in the cache.
	GetCredentialVersion(id string) (int64, error)
	// DeleteCredentialVersion removes the cached version number for a credential.
	DeleteCredentialVersion(id string) error
}
