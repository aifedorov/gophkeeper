// Package auth provides authentication and user management services for the GophKeeper server.
//
// This package implements the core authentication business logic including user registration,
// login, password hashing, encryption key derivation, and context management for authenticated requests.
package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/pkg/validator"
	"go.uber.org/zap"
)

// ContextKey is a type for context keys used to store user data in request context.
type ContextKey string

const (
	// userIDKey is the context key for storing user ID in request context.
	userIDKey = ContextKey("user_id")
	// encryptionKeyKey is the context key for storing base64-encoded encryption key in request context.
	encryptionKeyKey = ContextKey("encryption_key_encoded")
)

// Service defines the interface for authentication and user management operations.
// It provides methods for user registration, login, and context management for authenticated requests.
type Service interface {
	// Register creates a new auth account with the provided login and password.
	// Returns the user entity, the encryption key, and an error if the registration fails.
	Register(ctx context.Context, login, password string) (*User, []byte, error)
	// Login authenticates a auth with the provided credentials.
	// Returns the user entity, the encryption key, and an error if the login fails.
	Login(ctx context.Context, login, password string) (*User, []byte, error)
	// SetUserID Set userID in context.
	SetUserID(ctx context.Context, userID string) context.Context
	// GetUserIDFromContext List userID from context.
	GetUserIDFromContext(ctx context.Context) (string, error)
	// SetEncryptionKeyEncoded Set encryption key in base64 in context.
	SetEncryptionKeyEncoded(ctx context.Context, encryptionKey string) context.Context
	// GetEncryptionKeyFromContext List an encryption key in base64 from context.
	GetEncryptionKeyFromContext(ctx context.Context) (string, error)
	// GetUserDataFromContext retrieves both user ID and encryption key from the request context.
	// Returns an error if either value is missing or invalid.
	GetUserDataFromContext(ctx context.Context) (userID, encryptionKey string, err error)
}

// service implements the Service interface for authentication operations.
type service struct {
	repo      interfaces.Repository
	logger    *zap.Logger
	cryptoSrv interfaces.CryptoService
}

// NewService creates a new instance of the authentication service with the provided dependencies.
// It initializes the service with a user repository, logger, and cryptographic service.
func NewService(repo interfaces.Repository, logger *zap.Logger, cryptoSrv interfaces.CryptoService) Service {
	return &service{
		repo:      repo,
		logger:    logger,
		cryptoSrv: cryptoSrv,
	}
}

// Register creates a new user account with the provided login and password.
// It validates credentials, generates a salt, derives an encryption key, hashes the password,
// and stores the user in the repository. Returns the created user entity, encryption key, and any error.
// Returns ErrLoginExists if the login is already taken.
func (s *service) Register(ctx context.Context, login, password string) (*User, []byte, error) {
	s.logger.Debug("auth: starting registration", zap.String("login", login))

	if err := validator.ValidateLogin(login); err != nil {
		s.logger.Debug("auth: login validation failed", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: invalid login: %w", err)
	}
	if err := validator.ValidatePassword(password); err != nil {
		s.logger.Debug("auth: password validation failed", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: invalid password: %w", err)
	}

	salt, err := s.cryptoSrv.GenerateSalt()
	if err != nil {
		s.logger.Error("auth: failed to generate salt", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: failed to generate salt: %w", err)
	}
	s.logger.Debug("auth: salt generated successfully")

	encryptionKey := s.cryptoSrv.DeriveEncryptionKey(password, string(salt))
	s.logger.Debug("auth: encryption key derived")

	passwordHash, err := s.cryptoSrv.HashPassword(password)
	if err != nil {
		s.logger.Error("auth: failed to hash password", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: failed to hash password: %w", err)
	}
	s.logger.Debug("auth: password hashed successfully")

	usr, err := NewUser(login, string(salt))
	if err != nil {
		s.logger.Error("auth: failed to create user entity", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: failed to create user: %w", err)
	}

	repositoryUser := toRepositoryUser(usr, passwordHash)

	user, err := s.repo.CreateUser(ctx, repositoryUser, passwordHash)
	if errors.Is(err, ErrLoginExists) {
		s.logger.Debug("auth: login already exists", zap.String("login", login))
		return nil, nil, ErrLoginExists
	}
	if err != nil {
		s.logger.Error("auth: failed to create user in repository", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: failed to create user: %w", err)
	}
	s.logger.Debug("auth: user created in repository", zap.String("user_id", user.ID))

	dusr, err := toDomainUser(user)
	if err != nil {
		s.logger.Error("auth: failed to convert to domain user", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: failed to convert user to domain user: %w", err)
	}

	s.logger.Debug("auth: registration completed successfully", zap.String("user_id", dusr.GetUserID()))
	return &dusr, encryptionKey, nil
}

// Login authenticates a user with the provided credentials.
// It validates credentials, retrieves the user from the repository, verifies the password,
// and derives the encryption key. Returns the user entity, encryption key, and any error.
// Returns ErrUserNotFound if the user doesn't exist or ErrInvalidCredentials if the password is incorrect.
func (s *service) Login(ctx context.Context, login, password string) (*User, []byte, error) {
	s.logger.Debug("auth: starting login", zap.String("login", login))

	if err := validator.ValidateLogin(login); err != nil {
		s.logger.Debug("auth: login validation failed", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: invalid login: %w", err)
	}
	if err := validator.ValidatePassword(password); err != nil {
		s.logger.Debug("auth: password validation failed", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: invalid password: %w", err)
	}

	user, err := s.repo.GetUser(ctx, login)
	if errors.Is(err, ErrUserNotFound) {
		s.logger.Debug("auth: user not found", zap.String("login", login))
		return nil, nil, ErrUserNotFound
	}
	if err != nil {
		s.logger.Error("auth: failed to get user from repository", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: failed to get user: %w", err)
	}
	s.logger.Debug("auth: user found", zap.String("user_id", user.ID))

	encryptionKey := s.cryptoSrv.DeriveEncryptionKey(password, user.Salt)
	s.logger.Debug("auth: encryption key derived")

	err = s.cryptoSrv.CompareHashAndPassword(user.PasswordHash, password)
	if err != nil {
		s.logger.Debug("auth: password comparison failed", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: failed to compare hash and password: %w", err)
	}
	s.logger.Debug("auth: password verified successfully")

	dusr, err := toDomainUser(user)
	if err != nil {
		s.logger.Error("auth: failed to convert to domain user", zap.Error(err))
		return nil, nil, fmt.Errorf("auth: failed to convert user to domain user: %w", err)
	}

	s.logger.Debug("auth: login completed successfully", zap.String("user_id", dusr.GetUserID()))
	return &dusr, encryptionKey, nil
}

// GetUserIDFromContext retrieves the user ID from the request context.
// Returns an error if the user ID is not found or is invalid.
func (s *service) GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok || userID == "" {
		s.logger.Error("auth: user id not found in context")
		return "", errors.New("auth: failed to get user id from context")
	}
	return userID, nil
}

// GetEncryptionKeyFromContext retrieves the base64-encoded encryption key from the request context.
// Returns an error if the encryption key is not found or is invalid.
func (s *service) GetEncryptionKeyFromContext(ctx context.Context) (string, error) {
	encryptionKey, ok := ctx.Value(encryptionKeyKey).(string)
	if !ok || encryptionKey == "" {
		s.logger.Error("auth: encryption key not found in context")
		return "", errors.New("auth: failed to get encryption key from context")
	}
	return encryptionKey, nil
}

// SetUserID stores the user ID in the request context and returns a new context with the value.
// This is used by authentication interceptors to pass user information to service methods.
func (s *service) SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// SetEncryptionKeyEncoded stores the base64-encoded encryption key in the request context
// and returns a new context with the value. This is used by authentication interceptors
// to pass encryption key information to service methods.
func (s *service) SetEncryptionKeyEncoded(ctx context.Context, encryptionKey string) context.Context {
	return context.WithValue(ctx, encryptionKeyKey, encryptionKey)
}

// GetUserDataFromContext retrieves both user ID and encryption key from the request context.
// This is a convenience method that calls both GetUserIDFromContext and GetEncryptionKeyFromContext.
// Returns an error if either value is missing or invalid.
func (s *service) GetUserDataFromContext(ctx context.Context) (userID, encryptionKey string, err error) {
	userID, err = s.GetUserIDFromContext(ctx)
	if err != nil {
		return "", "", fmt.Errorf("grpc: failed to get userID: %w", err)
	}

	encryptionKey, err = s.GetEncryptionKeyFromContext(ctx)
	if err != nil {
		return "", "", fmt.Errorf("grpc: failed to get encryption key: %w", err)
	}

	return userID, encryptionKey, nil
}
