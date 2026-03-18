// Package credential provides credential management services for the GophKeeper server.
//
// This package implements the core business logic for managing user credentials (login/password pairs)
// with end-to-end encryption. All operations require user authentication and use encryption keys
// derived from the user's password.
package credential

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
	"go.uber.org/zap"
)

// Service defines the interface for credential management operations.
// All methods require user authentication and encryption key for data encryption/decryption.
type Service interface {
	// Create stores a new credential for the specified user with encryption.
	// The encryptionKey should be provided as a base64-encoded string.
	// Returns ErrNameExists if a credential with the same name already exists for the user.
	Create(ctx context.Context, userID, encryptionKey string, credential Credential) (*Credential, error)
	// List retrieves all credentials for the specified user and decrypts them.
	// The encryptionKey should be provided as a base64-encoded string.
	// Returns an empty slice if the user has no credentials.
	List(ctx context.Context, userID, encryptionKey string) ([]Credential, error)
	// Update modifies an existing credential for the specified user.
	// The encryptionKey should be provided as a base64-encoded string.
	// Returns ErrNotFound if the credential doesn't exist or doesn't belong to the user.
	Update(ctx context.Context, userID, encryptionKey string, credential Credential) (*Credential, error)
	// Delete removes a credential for the specified user.
	// Returns ErrNotFound if the credential doesn't exist or doesn't belong to the user.
	Delete(ctx context.Context, userID, id string) error
}

// service implements the Service interface for credential management.
type service struct {
	repo   interfaces.Repository
	crypto interfaces.CryptoService
	logger *zap.Logger
}

// NewService creates a new instance of the credential service with the provided dependencies.
// It initializes the service with a credential repository, cryptographic service, and logger.
func NewService(repo interfaces.Repository, crypto interfaces.CryptoService, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		crypto: crypto,
		logger: logger,
	}
}

// Create stores a new credential for the specified user with encryption.
// The encryptionKey should be provided as a base64-encoded string.
// Returns ErrNameExists if a credential with the same name already exists for the user.
func (s *service) Create(ctx context.Context, userID, encryptionKey string, credential Credential) (*Credential, error) {
	s.logger.Debug("credential: creating credential",
		zap.String("user_id", userID),
		zap.String("name", credential.GetName()))

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		s.logger.Error("credential: failed to decode encryption key", zap.Error(err))
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	rCredential, err := toRepositoryCredential(s.crypto, key, credential)
	if err != nil {
		s.logger.Error("credential: failed to convert to repository credential", zap.Error(err))
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}
	s.logger.Debug("credential: credential encrypted successfully")

	result, err := s.repo.CreateCredential(ctx, userID, rCredential)
	if errors.Is(err, ErrNameExists) {
		s.logger.Debug("credential: name already exists", zap.String("name", credential.GetName()))
		return nil, ErrNameExists
	}
	if err != nil {
		s.logger.Error("credential: failed to create in repository", zap.Error(err))
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}
	if result == nil {
		s.logger.Error("credential: repository returned nil")
		return nil, fmt.Errorf("failed to create credential: credential is nil")
	}
	s.logger.Debug("credential: created in repository", zap.String("id", result.ID))

	domainCred, err := toDomainCredential(s.crypto, key, *result)
	if err != nil {
		s.logger.Error("credential: failed to convert to domain credential", zap.Error(err))
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}

	s.logger.Debug("credential: created successfully", zap.String("id", domainCred.GetID()))
	return &domainCred, nil
}

// List retrieves all credentials for the specified user and decrypts them.
// The encryptionKey should be provided as a base64-encoded string.
// Returns an empty slice if the user has no credentials.
func (s *service) List(ctx context.Context, userID, encryptionKey string) ([]Credential, error) {
	s.logger.Debug("credential: listing credentials", zap.String("user_id", userID))

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		s.logger.Error("credential: failed to decode encryption key", zap.Error(err))
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	credentials, err := s.repo.ListCredentials(ctx, userID)
	if err != nil {
		s.logger.Error("credential: failed to list from repository", zap.Error(err))
		return nil, fmt.Errorf("failed to get list of credentials: %w", err)
	}
	s.logger.Debug("credential: retrieved from repository", zap.Int("count", len(credentials)))

	res := make([]Credential, len(credentials))
	for i, cred := range credentials {
		domainCred, err := toDomainCredential(s.crypto, key, cred)
		if err != nil {
			s.logger.Error("credential: failed to decrypt credential",
				zap.String("id", cred.ID),
				zap.Error(err))
			return nil, fmt.Errorf("failed to convert credential: %w", err)
		}
		res[i] = domainCred
	}

	s.logger.Debug("credential: list completed successfully", zap.Int("count", len(res)))
	return res, nil
}

// Update modifies an existing credential for the specified user.
// The encryptionKey should be provided as a base64-encoded string.
// Returns ErrNotFound if the credential doesn't exist or doesn't belong to the user.
func (s *service) Update(ctx context.Context, userID, encryptionKey string, credential Credential) (*Credential, error) {
	s.logger.Debug("credential: updating credential",
		zap.String("user_id", userID),
		zap.String("id", credential.GetID()))

	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		s.logger.Error("credential: failed to decode encryption key", zap.Error(err))
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	rCredential, err := toRepositoryCredential(s.crypto, key, credential)
	if err != nil {
		s.logger.Error("credential: failed to encrypt", zap.Error(err))
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}

	result, err := s.repo.UpdateCredential(ctx, userID, rCredential)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("credential: not found for update", zap.String("id", credential.GetID()))
		return nil, ErrNotFound
	}
	if errors.Is(err, ErrNameExists) {
		s.logger.Debug("credential: name already exists", zap.String("name", credential.GetName()))
		return nil, ErrNameExists
	}
	if errors.Is(err, ErrVersionConflict) {
		s.logger.Debug("credential: version conflict", zap.String("id", credential.GetID()))
		return nil, ErrVersionConflict
	}
	if err != nil || result == nil {
		s.logger.Error("credential: failed to update in repository", zap.Error(err))
		return nil, fmt.Errorf("failed to update credential: %w", err)
	}

	domainCred, err := toDomainCredential(s.crypto, key, *result)
	if err != nil {
		s.logger.Error("credential: failed to decrypt updated credential", zap.Error(err))
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}

	s.logger.Debug("credential: updated successfully", zap.String("id", domainCred.GetID()))
	return &domainCred, nil
}

// Delete removes a credential for the specified user.
// Returns ErrNotFound if the credential doesn't exist or doesn't belong to the user.
func (s *service) Delete(ctx context.Context, userID, id string) error {
	s.logger.Debug("credential: deleting credential",
		zap.String("user_id", userID),
		zap.String("id", id),
	)

	err := s.repo.DeleteCredential(ctx, userID, id)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("credential: not found for deletion", zap.String("id", id))
		return ErrNotFound
	}
	if err != nil {
		s.logger.Error("credential: failed to delete from repository", zap.Error(err))
		return fmt.Errorf("failed to delete credential: %w", err)
	}

	s.logger.Debug("credential: deleted successfully", zap.String("id", id))
	return nil
}
