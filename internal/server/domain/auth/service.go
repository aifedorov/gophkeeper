package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/pkg/validator"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContextKey string

const idKey = ContextKey("user_id")

type Service interface {
	// Register creates a new auth account with the provided login and password.
	// The password is hashed before storage. Returns ErrLoginExists if the login is already taken.
	Register(ctx context.Context, login, password string) (*User, error)
	// Login authenticates a auth with the provided credentials.
	// Returns ErrUserNotFound if the auth doesn't exist or if the password is incorrect.
	Login(ctx context.Context, login, password string) (*User, error)
	// GetUserIDFromContext returns the user id from the request context.
	GetUserIDFromContext(ctx context.Context) (string, error)
	// SetUserIDToContext sets the user id in the request context.
	SetUserIDToContext(ctx context.Context, userID string) context.Context
}

type service struct {
	repo         interfaces.Repository
	logger       *zap.Logger
	sessionStore SessionStore
	cryptoSrv    interfaces.CryptoService
}

func NewService(repo interfaces.Repository, logger *zap.Logger, cryptoSrv interfaces.CryptoService) Service {
	return &service{
		repo:         repo,
		logger:       logger,
		sessionStore: NewSessionStore(logger),
		cryptoSrv:    cryptoSrv,
	}
}

func (s *service) Register(ctx context.Context, login, password string) (*User, error) {
	s.logger.Debug("auth: starting registration", zap.String("login", login))

	if err := validator.ValidateLogin(login); err != nil {
		s.logger.Debug("auth: login validation failed", zap.Error(err))
		return nil, fmt.Errorf("auth: invalid login: %w", err)
	}
	if err := validator.ValidatePassword(password); err != nil {
		s.logger.Debug("auth: password validation failed", zap.Error(err))
		return nil, fmt.Errorf("auth: invalid password: %w", err)
	}

	salt, err := s.cryptoSrv.GenerateSalt()
	if err != nil {
		s.logger.Error("auth: failed to generate salt", zap.Error(err))
		return nil, fmt.Errorf("auth: failed to generate salt: %w", err)
	}
	s.logger.Debug("auth: salt generated successfully")

	encryptionKey := s.cryptoSrv.DeriveEncryptionKey(password, string(salt))
	s.logger.Debug("auth: encryption key derived")

	passwordHash, err := s.cryptoSrv.HashPassword(password)
	if err != nil {
		s.logger.Error("auth: failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("auth: failed to hash password: %w", err)
	}
	s.logger.Debug("auth: password hashed successfully")

	usr, err := NewUser(login, string(salt))
	if err != nil {
		s.logger.Error("auth: failed to create user entity", zap.Error(err))
		return nil, fmt.Errorf("auth: failed to create user: %w", err)
	}

	repositoryUser := toRepositoryUser(usr, passwordHash)

	user, err := s.repo.CreateUser(ctx, repositoryUser, passwordHash)
	if errors.Is(err, ErrLoginExists) {
		s.logger.Debug("auth: login already exists", zap.String("login", login))
		return nil, ErrLoginExists
	}
	if err != nil {
		s.logger.Error("auth: failed to create user in repository", zap.Error(err))
		return nil, fmt.Errorf("auth: failed to create user: %w", err)
	}
	s.logger.Debug("auth: user created in repository", zap.String("user_id", user.ID))

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		s.logger.Error("auth: failed to parse user id", zap.String("user_id", user.ID), zap.Error(err))
		return nil, fmt.Errorf("auth: failed to parse user id: %w", err)
	}

	s.sessionStore.Set(userID, encryptionKey)
	s.logger.Debug("auth: encryption key stored in session", zap.String("user_id", userID.String()))

	dusr, err := toDomainUser(user)
	if err != nil {
		s.logger.Error("auth: failed to convert to domain user", zap.Error(err))
		return nil, fmt.Errorf("auth: failed to convert user to domain user: %w", err)
	}

	s.logger.Debug("auth: registration completed successfully", zap.String("user_id", userID.String()))
	return &dusr, nil
}

func (s *service) Login(ctx context.Context, login, password string) (*User, error) {
	s.logger.Debug("auth: starting login", zap.String("login", login))

	if err := validator.ValidateLogin(login); err != nil {
		s.logger.Debug("auth: login validation failed", zap.Error(err))
		return nil, fmt.Errorf("auth: invalid login: %w", err)
	}
	if err := validator.ValidatePassword(password); err != nil {
		s.logger.Debug("auth: password validation failed", zap.Error(err))
		return nil, fmt.Errorf("auth: invalid password: %w", err)
	}

	user, err := s.repo.GetUser(ctx, login)
	if errors.Is(err, ErrUserNotFound) {
		s.logger.Debug("auth: user not found", zap.String("login", login))
		return nil, ErrUserNotFound
	}
	if err != nil {
		s.logger.Error("auth: failed to get user from repository", zap.Error(err))
		return nil, fmt.Errorf("auth: failed to get user: %w", err)
	}
	s.logger.Debug("auth: user found", zap.String("user_id", user.ID))

	err = s.cryptoSrv.CompareHashAndPassword(user.PasswordHash, password)
	if err != nil {
		s.logger.Debug("auth: password comparison failed", zap.Error(err))
		return nil, fmt.Errorf("auth: failed to compare hash and password: %w", err)
	}
	s.logger.Debug("auth: password verified successfully")

	dusr, err := toDomainUser(user)
	if err != nil {
		s.logger.Error("auth: failed to convert to domain user", zap.Error(err))
		return nil, fmt.Errorf("auth: failed to convert user to domain user: %w", err)
	}

	encryptionKey := s.cryptoSrv.DeriveEncryptionKey(password, dusr.GetSalt())
	s.logger.Debug("auth: encryption key derived")

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		s.logger.Error("auth: failed to parse user id", zap.String("user_id", user.ID), zap.Error(err))
		return nil, fmt.Errorf("auth: failed to parse user id: %w", err)
	}

	s.sessionStore.Set(userID, encryptionKey)
	s.logger.Debug("auth: encryption key stored in session", zap.String("user_id", userID.String()))

	s.logger.Debug("auth: login completed successfully", zap.String("user_id", userID.String()))
	return &dusr, nil
}

func (s *service) GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(idKey).(string)
	if !ok || userID == "" {
		s.logger.Debug("auth: user id not found in context")
		return "", errors.New("auth: failed to get user id from context")
	}
	s.logger.Debug("auth: user id extracted from context", zap.String("user_id", userID))
	return userID, nil
}

func (s *service) SetUserIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, idKey, userID)
}
