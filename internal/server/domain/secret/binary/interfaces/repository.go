package interfaces

import (
	"context"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock_repository.go -package=mocks

type Repository interface {
	CreateFile(ctx context.Context, userID string, file RepositoryFile) error
	ListFiles(ctx context.Context, userID string) ([]RepositoryFile, error)
	DeleteFile(ctx context.Context, userID, id string) error
}
