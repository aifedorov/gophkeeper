package credential

import (
	"context"
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error)
	Get(ctx context.Context, userID, id uuid.UUID) (*Credential, error)
	List(ctx context.Context, userID uuid.UUID) ([]Credential, error)
	Update(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error)
	Delete(ctx context.Context, userID, id uuid.UUID) error
}

type service struct {
	repo         interfaces.Repository
	crypto       interfaces.CryptoService
	sessionStore auth.SessionStore
	logger       *zap.Logger
}

func NewService(repo interfaces.Repository, crypto interfaces.CryptoService, sessionStore auth.SessionStore, logger *zap.Logger) Service {
	return &service{
		repo:         repo,
		crypto:       crypto,
		sessionStore: sessionStore,
		logger:       logger,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error) {
	s.logger.Debug("credential: creating credential",
		zap.String("user_id", userID.String()),
		zap.String("name", credential.GetName()))

	key, ok := s.sessionStore.GetEncryptionKey(userID)
	if !ok {
		s.logger.Debug("credential: encryption key not found in session", zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to get encryption key for user: %w", ErrNotFound)
	}
	s.logger.Debug("credential: encryption key retrieved from session")

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

	s.logger.Debug("credential: created successfully", zap.String("id", domainCred.GetID().String()))
	return &domainCred, nil
}

func (s *service) Get(ctx context.Context, userID, id uuid.UUID) (*Credential, error) {
	s.logger.Debug("credential: getting credential",
		zap.String("user_id", userID.String()),
		zap.String("id", id.String()))

	key, ok := s.sessionStore.GetEncryptionKey(userID)
	if !ok {
		s.logger.Debug("credential: encryption key not found in session", zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to get encryption key for user")
	}

	rCredential, err := s.repo.GetCredential(ctx, userID, id)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("credential: not found", zap.String("id", id.String()))
		return nil, ErrNotFound
	}
	if err != nil {
		s.logger.Error("credential: failed to get from repository", zap.Error(err))
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}
	if rCredential == nil {
		s.logger.Error("credential: repository returned nil")
		return nil, fmt.Errorf("failed to get credential: credential is nil")
	}

	domainCred, err := toDomainCredential(s.crypto, key, *rCredential)
	if err != nil {
		s.logger.Error("credential: failed to decrypt", zap.Error(err))
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}

	s.logger.Debug("credential: retrieved successfully", zap.String("id", id.String()))
	return &domainCred, nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]Credential, error) {
	s.logger.Debug("credential: listing credentials", zap.String("user_id", userID.String()))

	key, ok := s.sessionStore.GetEncryptionKey(userID)
	if !ok {
		s.logger.Debug("credential: encryption key not found in session", zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to get encryption key for user")
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

func (s *service) Update(ctx context.Context, userID uuid.UUID, credential Credential) (*Credential, error) {
	s.logger.Debug("credential: updating credential",
		zap.String("user_id", userID.String()),
		zap.String("id", credential.GetID().String()))

	key, ok := s.sessionStore.GetEncryptionKey(userID)
	if !ok {
		s.logger.Debug("credential: encryption key not found in session", zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to get encryption key for user")
	}

	rCredential, err := toRepositoryCredential(s.crypto, key, credential)
	if err != nil {
		s.logger.Error("credential: failed to encrypt", zap.Error(err))
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}

	result, err := s.repo.UpdateCredential(ctx, userID, rCredential)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("credential: not found for update", zap.String("id", credential.GetID().String()))
		return nil, ErrNotFound
	}
	if err != nil {
		s.logger.Error("credential: failed to update in repository", zap.Error(err))
		return nil, fmt.Errorf("failed to update credential: %w", err)
	}
	if result == nil {
		s.logger.Error("credential: repository returned nil")
		return nil, fmt.Errorf("failed to update credential: credential is nil")
	}

	domainCred, err := toDomainCredential(s.crypto, key, *result)
	if err != nil {
		s.logger.Error("credential: failed to decrypt updated credential", zap.Error(err))
		return nil, fmt.Errorf("failed to convert credential: %w", err)
	}

	s.logger.Debug("credential: updated successfully", zap.String("id", domainCred.GetID().String()))
	return &domainCred, nil
}

func (s *service) Delete(ctx context.Context, userID, id uuid.UUID) error {
	s.logger.Debug("credential: deleting credential",
		zap.String("user_id", userID.String()),
		zap.String("id", id.String()))

	err := s.repo.DeleteCredential(ctx, userID, id)
	if errors.Is(err, ErrNotFound) {
		s.logger.Debug("credential: not found for deletion", zap.String("id", id.String()))
		return ErrNotFound
	}
	if err != nil {
		s.logger.Error("credential: failed to delete from repository", zap.Error(err))
		return fmt.Errorf("failed to delete credential: %w", err)
	}

	s.logger.Debug("credential: deleted successfully", zap.String("id", id.String()))
	return nil
}
