package interfaces

import (
	"context"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock_repository.go -package=mocks

// RepositoryUser represents user data as stored in the repository.
type RepositoryUser struct {
	ID           string
	Login        string
	PasswordHash string
	Salt         string
}

// Repository defines the interface for user repository operations.
type Repository interface {
	// CreateUser creates a new user in the repository.
	// Returns ErrLoginExists if the login already exists.
	CreateUser(ctx context.Context, user RepositoryUser, passwordHash string) (RepositoryUser, error)
	// GetUser retrieves a user by login from the repository.
	// Returns ErrUserNotFound if the user doesn't exist.
	GetUser(ctx context.Context, login string) (RepositoryUser, error)
}
