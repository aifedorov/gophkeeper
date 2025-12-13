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

func (r *repository) CreateCredential(ctx context.Context, userID string, credential interfaces.RepositoryCredential) (*interfaces.RepositoryCredential, error) {
	r.logger.Debug("repo: creating credential",
		zap.String("user_id", userID),
		zap.String("name", credential.Name))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	dbCredential, err := r.queries.CreateCredential(ctx, CreateCredentialParams{
		UserID:            userUUID,
		Name:              credential.Name,
		Encryptedlogin:    credential.EncryptedLogin,
		Encryptedpassword: credential.EncryptedPassword,
		Encryptednotes:    credential.EncryptedNotes,
	})
	if conflictError(err) {
		r.logger.Debug("repo: credential name already exists", zap.String("name", credential.Name))
		return nil, credentialDomain.ErrNameExists
	}
	if err != nil {
		r.logger.Error("repo: failed to create credential", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to create credential: %w", err)
	}
	result := toInterfacesCredential(dbCredential)
	r.logger.Debug("repo: credential created successfully", zap.String("id", result.ID))
	return &result, nil
}

func (r *repository) GetCredential(ctx context.Context, userID, id string) (*interfaces.RepositoryCredential, error) {
	r.logger.Debug("repo: getting credential",
		zap.String("user_id", userID),
		zap.String("id", id))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("repo: failed to parse credential id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse credential id: %w", err)
	}

	dbCredential, err := r.queries.GetCredential(ctx, GetCredentialParams{
		ID:     idUUID,
		UserID: userUUID,
	})
	if notFoundError(err) {
		r.logger.Debug("repo: credential not found", zap.String("id", id))
		return nil, credentialDomain.ErrNotFound
	}
	if err != nil {
		r.logger.Error("repo: failed to get credential", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to get credential: %w", err)
	}
	result := toInterfacesCredential(dbCredential)
	r.logger.Debug("repo: credential found", zap.String("id", result.ID))
	return &result, nil
}

func (r *repository) ListCredentials(ctx context.Context, userID string) ([]interfaces.RepositoryCredential, error) {
	r.logger.Debug("repo: listing credentials", zap.String("user_id", userID))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	dbCredentials, err := r.queries.ListCredentials(ctx, userUUID)
	if notFoundError(err) {
		r.logger.Debug("repo: no credentials found", zap.String("user_id", userID))
		return []interfaces.RepositoryCredential{}, nil
	}
	if err != nil {
		r.logger.Error("repo: failed to list credentials", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to list credentials: %w", err)
	}
	result := make([]interfaces.RepositoryCredential, len(dbCredentials))
	for i, cred := range dbCredentials {
		result[i] = toInterfacesCredential(cred)
	}
	r.logger.Debug("repo: credentials listed successfully", zap.Int("count", len(result)))
	return result, nil
}

func (r *repository) UpdateCredential(ctx context.Context, userID string, credential interfaces.RepositoryCredential) (*interfaces.RepositoryCredential, error) {
	r.logger.Debug("repo: updating credential",
		zap.String("user_id", userID),
		zap.String("id", credential.ID))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	id, err := uuid.Parse(credential.ID)
	if err != nil {
		r.logger.Error("repo: failed to parse credential id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse credential id: %w", err)
	}

	dbCredential, err := r.queries.UpdateCredential(ctx, UpdateCredentialParams{
		ID:                id,
		UserID:            userUUID,
		Name:              credential.Name,
		Encryptedlogin:    credential.EncryptedLogin,
		Encryptedpassword: credential.EncryptedPassword,
		Encryptednotes:    credential.EncryptedNotes,
	})
	if notFoundError(err) {
		r.logger.Debug("repo: credential not found for update", zap.String("id", credential.ID))
		return nil, credentialDomain.ErrNotFound
	}
	if err != nil {
		r.logger.Error("repo: failed to update credential", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to update credential: %w", err)
	}
	result := toInterfacesCredential(dbCredential)
	r.logger.Debug("repo: credential updated successfully", zap.String("id", result.ID))
	return &result, nil
}

func (r *repository) DeleteCredential(ctx context.Context, userID, id string) error {
	r.logger.Debug("repo: deleting credential",
		zap.String("user_id", userID),
		zap.String("id", id))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("repo: failed to parse credential id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse credential id: %w", err)
	}

	err = r.queries.DeleteCredential(ctx, DeleteCredentialParams{
		ID:     idUUID,
		UserID: userUUID,
	})
	if notFoundError(err) {
		r.logger.Debug("repo: credential not found for deletion", zap.String("id", id))
		return credentialDomain.ErrNotFound
	}
	if err != nil {
		r.logger.Error("repo: failed to delete credential", zap.Error(err))
		return fmt.Errorf("repo: failed to delete credential: %w", err)
	}
	r.logger.Debug("repo: credential deleted successfully", zap.String("id", id))
	return nil
}

func toInterfacesCredential(credential Credential) interfaces.RepositoryCredential {
	return interfaces.RepositoryCredential{
		ID:                credential.ID.String(),
		UserID:            credential.UserID.String(),
		Name:              credential.Name,
		EncryptedLogin:    credential.Encryptedlogin,
		EncryptedPassword: credential.Encryptedpassword,
		EncryptedNotes:    credential.Encryptednotes,
	}
}
