// Package card provides card management services for the GophKeeper client.
//
// This package implements the client-side logic for managing user cards (payment card information).
// It communicates with the server via gRPC and handles card validation before sending requests.
package card

import (
	"context"
	"fmt"
)

// Service defines the interface for client-side card management operations.
type Service interface {
	// Create sends a request to create a new card on the server.
	// It validates the card before sending the request.
	// Returns an error if validation fails or if the server request fails.
	Create(ctx context.Context, card Card) error
	// List retrieves all cards for the authenticated user from the server.
	// Returns an empty slice if the user has no cards.
	List(ctx context.Context) ([]Card, error)
	// Update sends a request to update an existing card on the server.
	// It validates the card before sending the request.
	// Returns an error if validation fails, if the card doesn't exist, or if the server request fails.
	Update(ctx context.Context, id string, card Card) error
	// Delete removes a card by ID.
	// Returns an error if the card doesn't exist or if the deletion fails.
	Delete(ctx context.Context, id string) error
}

// service implements the Service interface for client-side card management.
type service struct {
	client Client
	cache  CacheStorage
}

// NewService creates a new instance of the card service with the provided gRPC client.
func NewService(client Client, cache CacheStorage) Service {
	return &service{
		client: client,
		cache:  cache,
	}
}

// Create sends a request to create a new card on the server.
// It validates the card before sending the request.
// Returns an error if validation fails or if the server request fails.
func (s *service) Create(ctx context.Context, card Card) error {
	if err := card.Validate(); err != nil {
		return fmt.Errorf("card: invalid card: %w", err)
	}

	id, version, err := s.client.Create(ctx, card)
	if err != nil {
		return fmt.Errorf("card: failed to create card: %w", err)
	}

	err = s.cache.SetCardVersion(id, version)
	if err != nil {
		return fmt.Errorf("card: failed to save card to cache: %w", err)
	}

	return nil
}

// List retrieves all cards for the authenticated user from the server.
// Returns an empty slice if the user has no cards.
func (s *service) List(ctx context.Context) ([]Card, error) {
	cards, err := s.client.List(ctx)
	if err != nil {
		return []Card{}, fmt.Errorf("card: failed to get list of cards: %w", err)
	}

	for _, card := range cards {
		if card.Version == 0 {
			return []Card{}, fmt.Errorf("card: server returned invalid version 0 for card %s", card.ID)
		}
		err := s.cache.SetCardVersion(card.ID, card.Version)
		if err != nil {
			return []Card{}, fmt.Errorf("card: failed to save card to cache: %w", err)
		}
	}

	return cards, nil
}

// Update sends a request to update an existing card on the server.
// It validates the card before sending the request.
// Returns an error if validation fails, if the card doesn't exist, or if the server request fails.
// Returns ErrVersionConflict if the card was modified by another client.
func (s *service) Update(ctx context.Context, id string, card Card) error {
	if err := card.Validate(); err != nil {
		return fmt.Errorf("card: invalid card: %w", err)
	}

	currentVersion, err := s.cache.GetCardVersion(id)
	if err != nil {
		return fmt.Errorf("card: failed to get version from cache (try running 'list' first): %w", err)
	}

	card.Version = currentVersion

	newVersion, err := s.client.Update(ctx, id, card)
	if err != nil {
		return propagateError("card: failed to update", err)
	}

	err = s.cache.SetCardVersion(card.ID, newVersion)
	if err != nil {
		return fmt.Errorf("card: failed to save card to cache: %w", err)
	}

	return nil
}

// Delete sends a request to delete a card by ID from the server.
// Returns an error if the card doesn't exist or if the deletion fails.
func (s *service) Delete(ctx context.Context, id string) error {
	err := s.client.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("card: failed to delete card: %w", err)
	}

	err = s.cache.DeleteCardVersion(id)
	if err != nil {
		return fmt.Errorf("card: failed to delete card from cache: %w", err)
	}

	return nil
}
