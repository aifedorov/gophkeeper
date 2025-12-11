package credential

import (
	"context"
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error)
	Get(ctx context.Context, userID, id uuid.UUID) (*Credential, error)
	List(ctx context.Context, userID uuid.UUID) ([]Credential, error)
	Update(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error)
	Delete(ctx context.Context, userID, id uuid.UUID) error
}

type service struct {
	repo interfaces.Repository
}

func NewService(repo interfaces.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error) {
	rCredential := toRepositoryCredential(credential)
	result, err := s.repo.CreateCredential(ctx, userID, rCredential)
	if errors.Is(err, ErrNameExists) {
		return nil, ErrNameExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("failed to create credential: credential is nil")
	}
	domainCred, err := toDomainCredential(*result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}
	return &domainCred, nil
}

func (s *service) Get(ctx context.Context, userID, id uuid.UUID) (*Credential, error) {
	rCredential, err := s.repo.GetCredential(ctx, userID, id)
	if errors.Is(err, ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}
	if rCredential == nil {
		return nil, fmt.Errorf("failed to get credential: credential is nil")
	}
	domainCred, err := toDomainCredential(*rCredential)
	if err != nil {
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}
	return &domainCred, nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]Credential, error) {
	credentials, err := s.repo.ListCredentials(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of credentials: %w", err)
	}

	res := make([]Credential, len(credentials))
	for i, cred := range credentials {
		domainCred, err := toDomainCredential(cred)
		if err != nil {
			return nil, fmt.Errorf("failed to convert credential: %w", err)
		}
		res[i] = domainCred
	}

	return res, nil
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error) {
	rCredential := toRepositoryCredential(credential)
	result, err := s.repo.UpdateCredential(ctx, userID, rCredential)
	if errors.Is(err, ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update credential: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("failed to update credential: credential is nil")
	}
	domainCred, err := toDomainCredential(*result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}
	return &domainCred, nil
}

func (s *service) Delete(ctx context.Context, userID, id uuid.UUID) error {
	err := s.repo.DeleteCredential(ctx, userID, id)
	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to delete credential: %w", err)
	}
	return nil
}
