package interfaces

import (
	"context"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock_repository.go -package=mocks

type RepositoryCard struct {
	ID                      string
	UserID                  string
	Name                    string
	EncryptedNumber         []byte
	EncryptedExpiredDate    []byte
	EncryptedCardHolderName []byte
	EncryptedCvv            []byte
	EncryptedNotes          []byte
}

type Repository interface {
	CreateCard(ctx context.Context, userID string, card RepositoryCard) (*RepositoryCard, error)
	ListCards(ctx context.Context, userID string) ([]RepositoryCard, error)
	UpdateCard(ctx context.Context, userID string, card RepositoryCard) (*RepositoryCard, error)
	DeleteCard(ctx context.Context, userID, id string) error
}
