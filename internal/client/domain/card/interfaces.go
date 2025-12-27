// Package card provides client interfaces for card management.
package card

import (
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=mock_interfaces_test.go -package=card

// Client defines the interface for gRPC client operations for card management.
// This interface abstracts the gRPC communication layer.
type Client interface {
	// Create sends a request to create a new card on the server.
	Create(ctx context.Context, card Card) (id string, version int64, err error)
	// Update sends a request to update an existing card on the server.
	// Returns the new version number after successful update.
	Update(ctx context.Context, id string, card Card) (version int64, err error)
	// Delete sends a request to delete a card by ID from the server.
	Delete(ctx context.Context, id string) error
	// List retrieves all cards for the authenticated user from the server.
	List(ctx context.Context) ([]Card, error)
}

// CacheStorage defines the interface for caching card version numbers locally.
// This is used for optimistic locking to detect concurrent modifications.
type CacheStorage interface {
	// SetCardVersion stores the version number for a card in the cache.
	SetCardVersion(id string, version int64) error
	// GetCardVersion retrieves the cached version number for a card.
	// Returns an error if the card is not found in the cache.
	GetCardVersion(id string) (int64, error)
	// DeleteCardVersion removes the cached version number for a card.
	DeleteCardVersion(id string) error
}
