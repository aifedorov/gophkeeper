package repository

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=querier.go -destination=mock_querier_test.go -package=repository

type Querier interface {
	CreateFile(ctx context.Context, arg CreateFileParams) error
	GetFile(ctx context.Context, arg GetFileParams) (File, error)
	GetFileForUpdate(ctx context.Context, arg GetFileForUpdateParams) (File, error)
	ListFiles(ctx context.Context, userID uuid.UUID) ([]File, error)
	DeleteFile(ctx context.Context, arg DeleteFileParams) (int64, error)
	UpdateFile(ctx context.Context, arg UpdateFileParams) error
}
