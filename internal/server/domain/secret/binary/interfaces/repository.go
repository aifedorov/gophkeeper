package interfaces

import (
	"context"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock_repository.go -package=mocks

type Repository interface {
	Create(ctx context.Context, userID string, file RepositoryFile) error
	Get(ctx context.Context, userID, id string) (RepositoryFile, error)
	List(ctx context.Context, userID string) ([]RepositoryFile, error)
	Delete(ctx context.Context, userID, id string) error
}
