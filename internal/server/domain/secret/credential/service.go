package credential

import (
	"context"
	"errors"
	"fmt"

	repository "github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/repository/db"
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
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error) {
	rCredential, err := s.repo.CreateCredential(ctx, userID, toRepositoryCredential(credential))
	if errors.Is(err, repository.ErrNameExists) {
		return nil, ErrNameExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}
	if rCredential == nil {
		return nil, fmt.Errorf("failed to create credential: credential is nil")
	}
	credential = toDomainCredential(*rCredential)
	return &credential, nil
}

func (s *service) Get(ctx context.Context, userID, id uuid.UUID) (*Credential, error) {
	rCredential, err := s.repo.GetCredential(ctx, userID, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}
	if rCredential == nil {
		return nil, fmt.Errorf("failed to get credential: credential is nil")
	}
	credential := toDomainCredential(*rCredential)
	return &credential, nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]Credential, error) {
	credentials, err := s.repo.ListCredentials(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of credentials: %w", err)
	}

	res := make([]Credential, len(credentials))
	for i, cred := range credentials {
		res[i] = toDomainCredential(cred)
	}

	return res, nil
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error) {
	rCredential, err := s.repo.UpdateCredential(ctx, userID, toRepositoryCredential(credential))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update credential: %w", err)
	}
	if rCredential == nil {
		return nil, fmt.Errorf("failed to update credential: credential is nil")
	}
	credential = toDomainCredential(*rCredential)
	return &credential, nil
}

func (s *service) Delete(ctx context.Context, userID, id uuid.UUID) error {
	err := s.repo.DeleteCredential(ctx, userID, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to delete credential: %w", err)
	}
	return nil
}
