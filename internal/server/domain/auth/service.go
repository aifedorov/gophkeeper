package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/pkg/validator"
	"go.uber.org/zap"
)

type ContextKey string

const (
	userIDKey        = ContextKey("user_id")
	encryptionKeyKey = ContextKey("encryption_key_encoded")
)

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
}

type service struct {
	repo      interfaces.Repository
	logger    *zap.Logger
	cryptoSrv interfaces.CryptoService
}

func NewService(repo interfaces.Repository, logger *zap.Logger, cryptoSrv interfaces.CryptoService) Service {
	return &service{
		repo:      repo,
		logger:    logger,
		cryptoSrv: cryptoSrv,
	}
}

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

func (s *service) GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok || userID == "" {
		s.logger.Error("auth: user id not found in context")
		return "", errors.New("auth: failed to get user id from context")
	}
	return ctx.Value(userIDKey).(string), nil
}

func (s *service) GetEncryptionKeyFromContext(ctx context.Context) (string, error) {
	return ctx.Value(encryptionKeyKey).(string), nil
}

func (s *service) SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func (s *service) SetEncryptionKeyEncoded(ctx context.Context, encryptionKey string) context.Context {
	return context.WithValue(ctx, encryptionKeyKey, encryptionKey)
}
