package repository

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=querier.go -destination=mock_querier_test.go -package=repository

type Querier interface {
	CreateCredential(ctx context.Context, arg CreateCredentialParams) (Credential, error)
	ListCredentials(ctx context.Context, userID uuid.UUID) ([]Credential, error)
	UpdateCredential(ctx context.Context, arg UpdateCredentialParams) (Credential, error)
	DeleteCredential(ctx context.Context, arg DeleteCredentialParams) error
}
