package repository

import (
	"context"
	"fmt"

	credentialDomain "github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type repository struct {
	queries Querier
	logger  *zap.Logger
}

func NewRepository(db DBTX, logger *zap.Logger) interfaces.Repository {
	return &repository{
		queries: New(db),
		logger:  logger,
	}
}

func NewRepositoryWithQuerier(querier Querier, logger *zap.Logger) interfaces.Repository {
	return &repository{
		queries: querier,
		logger:  logger,
	}
}

func (r *repository) CreateCredential(ctx context.Context, userID uuid.UUID, credential interfaces.RepositoryCredential) (*interfaces.RepositoryCredential, error) {
	dbCredential, err := r.queries.CreateCredential(ctx, CreateCredentialParams{
		UserID:            userID,
		Name:              credential.Name,
		Encryptedlogin:    credential.Encryptedlogin,
		Encryptedpassword: credential.Encryptedpassword,
		Encryptednotes:    credential.Encryptednotes,
	})
	if conflictError(err) {
		return nil, credentialDomain.ErrNameExists
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to create credential: %w", err)
	}
	result := toInterfacesCredential(dbCredential)
	return &result, nil
}

func (r *repository) GetCredential(ctx context.Context, userID, id uuid.UUID) (*interfaces.RepositoryCredential, error) {
	dbCredential, err := r.queries.GetCredential(ctx, GetCredentialParams{
		ID:     id,
		UserID: userID,
	})
	if notFoundError(err) {
		return nil, credentialDomain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to get credential: %w", err)
	}
	result := toInterfacesCredential(dbCredential)
	return &result, nil
}

func (r *repository) ListCredentials(ctx context.Context, userID uuid.UUID) ([]interfaces.RepositoryCredential, error) {
	dbCredentials, err := r.queries.ListCredentials(ctx, userID)
	if notFoundError(err) {
		return []interfaces.RepositoryCredential{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to list credentials: %w", err)
	}
	result := make([]interfaces.RepositoryCredential, len(dbCredentials))
	for i, cred := range dbCredentials {
		result[i] = toInterfacesCredential(cred)
	}
	return result, nil
}

func (r *repository) UpdateCredential(ctx context.Context, userID uuid.UUID, credential interfaces.RepositoryCredential) (*interfaces.RepositoryCredential, error) {
	id, err := uuid.Parse(credential.ID)
	if err != nil {
		return nil, fmt.Errorf("repo: failed to parse credential id: %w", err)
	}
	userIDParsed, err := uuid.Parse(credential.UserID)
	if err != nil {
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	dbCredential, err := r.queries.UpdateCredential(ctx, UpdateCredentialParams{
		ID:                id,
		UserID:            userIDParsed,
		Name:              credential.Name,
		Encryptedlogin:    credential.Encryptedlogin,
		Encryptedpassword: credential.Encryptedpassword,
		Encryptednotes:    credential.Encryptednotes,
	})
	if notFoundError(err) {
		return nil, credentialDomain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to update credential: %w", err)
	}
	result := toInterfacesCredential(dbCredential)
	return &result, nil
}

func (r *repository) DeleteCredential(ctx context.Context, userID, id uuid.UUID) error {
	err := r.queries.DeleteCredential(ctx, DeleteCredentialParams{
		ID:     id,
		UserID: userID,
	})
	if notFoundError(err) {
		return credentialDomain.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("repo: failed to delete credential: %w", err)
	}
	return nil
}

func toInterfacesCredential(credential Credential) interfaces.RepositoryCredential {
	return interfaces.RepositoryCredential{
		ID:                credential.ID.String(),
		UserID:            credential.UserID.String(),
		Name:              credential.Name,
		Encryptedlogin:    credential.Encryptedlogin,
		Encryptedpassword: credential.Encryptedpassword,
		Encryptednotes:    credential.Encryptednotes,
	}
}
