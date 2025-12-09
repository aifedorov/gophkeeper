package repository

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateCredential(ctx context.Context, arg CreateCredentialParams) (Credential, error)
	GetCredential(ctx context.Context, arg GetCredentialParams) (Credential, error)
	ListCredentials(ctx context.Context, userID uuid.UUID) ([]Credential, error)
	UpdateCredential(ctx context.Context, arg UpdateCredentialParams) (Credential, error)
	DeleteCredential(ctx context.Context, arg DeleteCredentialParams) error
}
