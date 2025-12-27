// Package card provides card management services for the GophKeeper server.
//
// This package implements the core business logic for managing payment card information
// with end-to-end encryption. All operations require user authentication and use encryption keys
// derived from the user's password.
package card

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/card/interfaces"
	"go.uber.org/zap"
)

// Service defines the interface for card management operations.
// All methods require user authentication and encryption key for data encryption/decryption.
type Service interface {
	// Create stores a new card for the specified user with encryption.
	// The encryptionKey should be provided as a base64-encoded string.
	// Returns ErrNameExists if a card with the same name already exists for the user.
	Create(ctx context.Context, userID, encryptionKey string, card Card) (*Card, error)
	// List retrieves all cards for the specified user and decrypts them.
	// The encryptionKey should be provided as a base64-encoded string.
	// Returns an empty slice if the user has no cards.
	List(ctx context.Context, userID, encryptionKey string) ([]Card, error)
	// Update modifies an existing card for the specified user.
	// The encryptionKey should be provided as a base64-encoded string.
	// Returns ErrNotFound if the card doesn't exist or doesn't belong to the user.
	Update(ctx context.Context, userID, encryptionKey string, card Card) (*Card, error)
	// Delete removes a card for the specified user.
	// Returns ErrNotFound if the card doesn't exist or doesn't belong to the user.
	Delete(ctx context.Context, userID, id string) error
}

// service implements the Service interface for card management.
type service struct {
	repo   interfaces.Repository
	crypto interfaces.CryptoService
	logger *zap.Logger
}

// NewService creates a new instance of the card service with the provided dependencies.
// It initializes the service with a card repository, cryptographic service, and logger.
func NewService(repo interfaces.Repository, crypto interfaces.CryptoService, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		crypto: crypto,
		logger: logger,
	}
}

// Create stores a new card for the specified user with encryption.
// The encryptionKey should be provided as a base64-encoded string.
// Returns ErrNameExists if a card with the same name already exists for the user.
func (s *service) Create(ctx context.Context, userID, encryptionKey string, card Card) (*Card, error) {
	s.logger.Debug("card: creating card",
		zap.String("user_id", userID),
		zap.String("name", card.GetName()))

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		s.logger.Error("card: failed to decode encryption key", zap.Error(err))
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	rCard, err := toRepositoryCard(s.crypto, key, card)
	if err != nil {
		s.logger.Error("card: failed to convert to repository card", zap.Error(err))
		return nil, fmt.Errorf("failed to convert card: %w", err)
	}
	s.logger.Debug("card: card encrypted successfully")

	result, err := s.repo.CreateCard(ctx, userID, rCard)
	if errors.Is(err, ErrNameExists) {
		s.logger.Debug("card: name already exists", zap.String("name", card.GetName()))
		return nil, ErrNameExists
	}
	if err != nil {
		s.logger.Error("card: failed to create in repository", zap.Error(err))
		return nil, fmt.Errorf("failed to create card: %w", err)
	}
	if result == nil {
		s.logger.Error("card: repository returned nil")
		return nil, fmt.Errorf("failed to create card: card is nil")
	}
	s.logger.Debug("card: created in repository", zap.String("id", result.ID))

	domainCard, err := toDomainCard(s.crypto, key, *result)
	if err != nil {
		s.logger.Error("card: failed to convert to domain card", zap.Error(err))
		return nil, fmt.Errorf("failed to convert card: %w", err)
	}

	s.logger.Debug("card: created successfully", zap.String("id", domainCard.GetID()))
	return &domainCard, nil
}

// List retrieves all cards for the specified user and decrypts them.
// The encryptionKey should be provided as a base64-encoded string.
// Returns an empty slice if the user has no cards.
func (s *service) List(ctx context.Context, userID, encryptionKey string) ([]Card, error) {
	s.logger.Debug("card: listing cards", zap.String("user_id", userID))

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		s.logger.Error("card: failed to decode encryption key", zap.Error(err))
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	cards, err := s.repo.ListCards(ctx, userID)
	if err != nil {
		s.logger.Error("card: failed to list from repository", zap.Error(err))
		return nil, fmt.Errorf("failed to get list of cards: %w", err)
	}
	s.logger.Debug("card: retrieved from repository", zap.Int("count", len(cards)))

	res := make([]Card, len(cards))
	for i, card := range cards {
		domainCard, err := toDomainCard(s.crypto, key, card)
		if err != nil {
			s.logger.Error("card: failed to decrypt card",
				zap.String("id", card.ID),
				zap.Error(err))
			return nil, fmt.Errorf("failed to convert card: %w", err)
		}
		res[i] = domainCard
	}

	s.logger.Debug("card: list completed successfully", zap.Int("count", len(res)))
	return res, nil
}

// Update modifies an existing card for the specified user.
// The encryptionKey should be provided as a base64-encoded string.
// Returns ErrNotFound if the card doesn't exist or doesn't belong to the user.
func (s *service) Update(ctx context.Context, userID, encryptionKey string, card Card) (*Card, error) {
	s.logger.Debug("card: updating card",
		zap.String("user_id", userID),
		zap.String("id", card.GetID()))

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		s.logger.Error("card: failed to decode encryption key", zap.Error(err))
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	rCard, err := toRepositoryCard(s.crypto, key, card)
	if err != nil {
		s.logger.Error("card: failed to encrypt", zap.Error(err))
		return nil, fmt.Errorf("failed to convert card: %w", err)
	}

	result, err := s.repo.UpdateCard(ctx, userID, rCard)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("card: not found for update", zap.String("id", card.GetID()))
		return nil, ErrNotFound
	}
	if errors.Is(err, ErrNameExists) {
		s.logger.Debug("card: name already exists", zap.String("name", card.GetName()))
		return nil, ErrNameExists
	}
	if errors.Is(err, ErrVersionConflict) {
		s.logger.Debug("card: version conflict", zap.String("id", card.GetID()))
		return nil, ErrVersionConflict
	}
	if err != nil || result == nil {
		s.logger.Error("card: failed to update in repository", zap.Error(err))
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	domainCard, err := toDomainCard(s.crypto, key, *result)
	if err != nil {
		s.logger.Error("card: failed to decrypt updated card", zap.Error(err))
		return nil, fmt.Errorf("failed to convert card: %w", err)
	}

	s.logger.Debug("card: updated successfully", zap.String("id", domainCard.GetID()))
	return &domainCard, nil
}

// Delete removes a card for the specified user.
// Returns ErrNotFound if the card doesn't exist or doesn't belong to the user.
func (s *service) Delete(ctx context.Context, userID, id string) error {
	s.logger.Debug("card: deleting card",
		zap.String("user_id", userID),
		zap.String("id", id),
	)

	err := s.repo.DeleteCard(ctx, userID, id)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("card: not found for deletion", zap.String("id", id))
		return ErrNotFound
	}
	if err != nil {
		s.logger.Error("card: failed to delete from repository", zap.Error(err))
		return fmt.Errorf("failed to delete card: %w", err)
	}

	s.logger.Debug("card: deleted successfully", zap.String("id", id))
	return nil
}
