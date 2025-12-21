// Package credential provides client interfaces for credential management.
package credential

import (
	"context"
)

// Client defines the interface for gRPC client operations for credential management.
// This interface abstracts the gRPC communication layer.
type Client interface {
	// Create sends a request to create a new credential on the server.
	Create(ctx context.Context, creds Credential) error
	// Update sends a request to update an existing credential on the server.
	Update(ctx context.Context, id string, creds Credential) error
	// Delete sends a request to delete a credential by ID from the server.
	Delete(ctx context.Context, id string) error
	// List retrieves all credentials for the authenticated user from the server.
	List(ctx context.Context) ([]Credential, error)
}
