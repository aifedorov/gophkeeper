package repository

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=querier.go -destination=mock_querier_test.go -package=repository

type Querier interface {
	CreateCard(ctx context.Context, arg CreateCardParams) (Card, error)
	ListCards(ctx context.Context, userID uuid.UUID) ([]Card, error)
	UpdateCard(ctx context.Context, arg UpdateCardParams) (Card, error)
	DeleteCard(ctx context.Context, arg DeleteCardParams) (int64, error)
}
