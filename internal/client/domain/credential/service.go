// Package credential provides credential management services for the GophKeeper client.
//
// This package implements the client-side logic for managing user credentials (login/password pairs).
// It communicates with the server via gRPC and handles credential validation before sending requests.
package credential

import (
	"context"
	"fmt"
)

// Service defines the interface for client-side credential management operations.
type Service interface {
	// Create sends a request to create a new credential on the server.
	// It validates the credential before sending the request.
	// Returns an error if validation fails or if the server request fails.
	Create(ctx context.Context, creds Credential) error
	// List retrieves all credentials for the authenticated user from the server.
	// Returns an empty slice if the user has no credentials.
	List(ctx context.Context) ([]Credential, error)
	// Update sends a request to update an existing credential on the server.
	// It validates the credential before sending the request.
	// Returns an error if validation fails, if the credential doesn't exist, or if the server request fails.
	Update(ctx context.Context, id string, cred Credential) error
	// Delete removes a credential by ID.
	// Returns an error if the credential doesn't exist or if the deletion fails.
	Delete(ctx context.Context, id string) error
}

// service implements the Service interface for client-side credential management.
type service struct {
	client Client
}

// NewService creates a new instance of the credential service with the provided gRPC client.
func NewService(client Client) Service {
	return &service{
		client: client,
	}
}

// Create sends a request to create a new credential on the server.
// It validates the credential before sending the request.
// Returns an error if validation fails or if the server request fails.
func (s *service) Create(ctx context.Context, creds Credential) error {
	if err := creds.Validate(); err != nil {
		return fmt.Errorf("credential: invalid credential: %w", err)
	}

	err := s.client.Create(ctx, creds)
	if err != nil {
		return fmt.Errorf("credential: failed to create credential: %w", err)
	}
	return nil
}

// List retrieves all credentials for the authenticated user from the server.
// Returns an empty slice if the user has no credentials.
func (s *service) List(ctx context.Context) ([]Credential, error) {
	creds, err := s.client.List(ctx)
	if err != nil {
		return []Credential{}, fmt.Errorf("credential: failed to get list of credentials: %w", err)
	}
	return creds, nil
}

// Update sends a request to update an existing credential on the server.
// It validates the credential before sending the request.
// Returns an error if validation fails, if the credential doesn't exist, or if the server request fails.
func (s *service) Update(ctx context.Context, id string, cred Credential) error {
	if err := cred.Validate(); err != nil {
		return fmt.Errorf("credential: invalid credential: %w", err)
	}

	err := s.client.Update(ctx, id, cred)
	if err != nil {
		return fmt.Errorf("credential: failed to get credential: %w", err)
	}
	return nil
}

// Delete sends a request to delete a credential by ID from the server.
// Returns an error if the credential doesn't exist or if the deletion fails.
func (s *service) Delete(ctx context.Context, id string) error {
	err := s.client.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("credential: failed to delete credential: %w", err)
	}
	return nil
}
