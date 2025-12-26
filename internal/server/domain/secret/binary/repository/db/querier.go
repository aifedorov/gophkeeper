package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=querier.go -destination=mock_querier_test.go -package=repository

type Querier interface {
	CreateFile(ctx context.Context, arg CreateFileParams) (File, error)
	GetFile(ctx context.Context, arg GetFileParams) (File, error)
	ListFiles(ctx context.Context, userID uuid.UUID) ([]File, error)
	GetFileForUpdate(ctx context.Context, arg GetFileForUpdateParams) (File, error)
	UpdateFile(ctx context.Context, arg UpdateFileParams) (File, error)
	DeleteFile(ctx context.Context, arg DeleteFileParams) (int64, error)
	WithTx(tx pgx.Tx) Querier
}
