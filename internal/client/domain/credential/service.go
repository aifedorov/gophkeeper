package credential

import (
	"context"
	"fmt"
)

type Service interface {
	Create(ctx context.Context, creds Credential) error
	List(ctx context.Context) ([]Credential, error)
	Update(ctx context.Context, id string, cred Credential) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	client Client
}

func NewService(client Client) Service {
	return &service{
		client: client,
	}
}

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

func (s *service) List(ctx context.Context) ([]Credential, error) {
	creds, err := s.client.List(ctx)
	if err != nil {
		return []Credential{}, fmt.Errorf("credential: failed to get list of credentials: %w", err)
	}
	return creds, nil
}

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

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.client.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("credential: failed to delete credential: %w", err)
	}
	return nil
}
