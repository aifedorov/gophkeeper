// Package card provides client interfaces for card management.
package card

import (
	"context"
)

// Client defines the interface for gRPC client operations for card management.
// This interface abstracts the gRPC communication layer.
type Client interface {
	// Create sends a request to create a new card on the server.
	Create(ctx context.Context, card Card) error
	// Update sends a request to update an existing card on the server.
	Update(ctx context.Context, id string, card Card) error
	// Delete sends a request to delete a card by ID from the server.
	Delete(ctx context.Context, id string) error
	// List retrieves all cards for the authenticated user from the server.
	List(ctx context.Context) ([]Card, error)
}
