package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	CreateCredential(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error)
	GetCredential(ctx context.Context, userID, id uuid.UUID) (*Credential, error)
	ListCredentials(ctx context.Context, userID uuid.UUID) ([]Credential, error)
	UpdateCredential(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error)
	DeleteCredential(ctx context.Context, userID, id uuid.UUID) error
}

type repository struct {
	queries Querier
	logger  *zap.Logger
}

func NewRepository(db DBTX, logger *zap.Logger) Repository {
	return &repository{
		queries: New(db),
		logger:  logger,
	}
}

func NewRepositoryWithQuerier(querier Querier, logger *zap.Logger) Repository {
	return &repository{
		queries: querier,
		logger:  logger,
	}
}

func (r *repository) CreateCredential(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error) {
	credential, err := r.queries.CreateCredential(ctx, CreateCredentialParams{
		UserID:   userID,
		Name:     credential.Name,
		Login:    credential.Login,
		Password: credential.Password,
		Metadata: credential.Metadata,
	})
	if conflictError(err) {
		return nil, ErrNameExists
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to create credential: %w", err)
	}
	return &credential, nil
}

func (r *repository) GetCredential(ctx context.Context, userID, id uuid.UUID) (*Credential, error) {
	credential, err := r.queries.GetCredential(ctx, GetCredentialParams{
		ID:     id,
		UserID: userID,
	})
	if notFoundError(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to get credential: %w", err)
	}
	return &credential, nil
}

func (r *repository) ListCredentials(ctx context.Context, userID uuid.UUID) ([]Credential, error) {
	credentials, err := r.queries.ListCredentials(ctx, userID)
	if notFoundError(err) {
		return []Credential{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to list credentials: %w", err)
	}
	return credentials, nil
}

func (r *repository) UpdateCredential(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error) {
	credential, err := r.queries.UpdateCredential(ctx, UpdateCredentialParams{
		ID:       credential.ID,
		UserID:   userID,
		Name:     credential.Name,
		Login:    credential.Login,
		Password: credential.Password,
		Metadata: credential.Metadata,
	})
	if notFoundError(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to update credential: %w", err)
	}
	return &credential, nil
}

func (r *repository) DeleteCredential(ctx context.Context, userID, id uuid.UUID) error {
	err := r.queries.DeleteCredential(ctx, DeleteCredentialParams{
		ID:     id,
		UserID: userID,
	})
	if notFoundError(err) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("repo: failed to delete credential: %w", err)
	}
	return nil
}
