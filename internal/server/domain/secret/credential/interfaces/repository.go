package interfaces

import (
	"context"

	"github.com/google/uuid"
)

// RepositoryCredential represents credential data as stored in the repository.
type RepositoryCredential struct {
	ID                string
	UserID            string
	Name              string
	Encryptedlogin    []byte
	Encryptedpassword []byte
	Encryptednotes    []byte
}

// Repository defines the interface for credential repository operations.
type Repository interface {
	// CreateCredential creates a new credential in the repository.
	// Returns ErrNameExists if a credential with the same name already exists for the user.
	CreateCredential(ctx context.Context, userID uuid.UUID, credential RepositoryCredential) (*RepositoryCredential, error)
	// GetCredential retrieves a credential by ID for the specified user.
	// Returns ErrNotFound if the credential doesn't exist.
	GetCredential(ctx context.Context, userID, id uuid.UUID) (*RepositoryCredential, error)
	// ListCredentials retrieves all credentials for the specified user.
	ListCredentials(ctx context.Context, userID uuid.UUID) ([]RepositoryCredential, error)
	// UpdateCredential updates an existing credential in the repository.
	// Returns ErrNotFound if the credential doesn't exist.
	UpdateCredential(ctx context.Context, userID uuid.UUID, credential RepositoryCredential) (*RepositoryCredential, error)
	// DeleteCredential soft deletes a credential by ID for the specified user.
	// Returns ErrNotFound if the credential doesn't exist.
	DeleteCredential(ctx context.Context, userID, id uuid.UUID) error
}
