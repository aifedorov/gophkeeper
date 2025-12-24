package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=querier.go -destination=mock_querier_test.go -package=repository

// Querier defines the interface for database query operations on credentials.
// This interface wraps sqlc-generated queries and adds transaction support.
type Querier interface {
	// CreateCredential inserts a new credential into the database.
	CreateCredential(ctx context.Context, arg CreateCredentialParams) (Credential, error)
	// ListCredentials retrieves all non-deleted credentials for a user.
	ListCredentials(ctx context.Context, userID uuid.UUID) ([]Credential, error)
	// GetCredentialForUpdate retrieves a credential with a row lock for update operations.
	GetCredentialForUpdate(ctx context.Context, arg GetCredentialForUpdateParams) (Credential, error)
	// UpdateCredential modifies an existing credential in the database.
	UpdateCredential(ctx context.Context, arg UpdateCredentialParams) (Credential, error)
	// DeleteCredential soft-deletes a credential by setting deleted_at timestamp.
	DeleteCredential(ctx context.Context, arg DeleteCredentialParams) (int64, error)
	// WithTx returns a new Querier that executes queries within the given transaction.
	WithTx(tx pgx.Tx) Querier
}
