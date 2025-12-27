package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=querier.go -destination=mock_querier_test.go -package=repository

// Querier defines the interface for database query operations on cards.
// This interface wraps sqlc-generated queries and adds transaction support.
type Querier interface {
	// CreateCard inserts a new card into the database.
	CreateCard(ctx context.Context, arg CreateCardParams) (Card, error)
	// ListCards retrieves all non-deleted cards for a user.
	ListCards(ctx context.Context, userID uuid.UUID) ([]Card, error)
	// GetCardForUpdate retrieves a card with a row lock for update operations.
	GetCardForUpdate(ctx context.Context, arg GetCardForUpdateParams) (Card, error)
	// UpdateCard modifies an existing card in the database.
	UpdateCard(ctx context.Context, arg UpdateCardParams) (Card, error)
	// DeleteCard soft-deletes a card by setting deleted_at timestamp.
	DeleteCard(ctx context.Context, arg DeleteCardParams) (int64, error)
	// WithTx returns a new Querier that executes queries within the given transaction.
	WithTx(tx pgx.Tx) Querier
}
