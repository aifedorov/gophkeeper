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
	sessionStore *SessionStore
	cryptoSrv    interfaces.CryptoService
}

func NewService(repo interfaces.Repository, logger *zap.Logger, cryptoSrv interfaces.CryptoService) Service {
	return &service{
		repo:         repo,
		logger:       logger,
		sessionStore: NewSessionStore(),
		cryptoSrv:    cryptoSrv,
	}
}

func (s *service) Register(ctx context.Context, login, password string) (*User, error) {
	if err := validator.ValidateLogin(login); err != nil {
		return nil, fmt.Errorf("invalid login: %w", err)
	}
	if err := validator.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	salt, err := s.cryptoSrv.GenerateSalt()
	if err != nil {
		s.logger.Error("failed to generate salt", zap.Error(err))
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	encryptionKey := s.cryptoSrv.DeriveEncryptionKey(password, string(salt))

	passwordHash, err := s.cryptoSrv.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	usr, err := NewUser(login, string(salt))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	repositoryUser := toRepositoryUser(usr, passwordHash)

	user, err := s.repo.CreateUser(ctx, repositoryUser, passwordHash)
	if errors.Is(err, ErrLoginExists) {
		return nil, ErrLoginExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %w", err)
	}
	s.sessionStore.Set(userID, encryptionKey)

	dusr, err := toDomainUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user to domain user: %w", err)
	}

	return &dusr, nil
}

func (s *service) Login(ctx context.Context, login, password string) (*User, error) {
	if err := validator.ValidateLogin(login); err != nil {
		return nil, fmt.Errorf("invalid login: %w", err)
	}
	if err := validator.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	user, err := s.repo.GetUser(ctx, login)
	if errors.Is(err, ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	err = s.cryptoSrv.CompareHashAndPassword(user.PasswordHash, password)
	if err != nil {
		return nil, fmt.Errorf("failed to compare hash and password: %w", err)
	}

	dusr, err := toDomainUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user to domain user: %w", err)
	}

	encryptionKey := s.cryptoSrv.DeriveEncryptionKey(password, dusr.GetSalt())
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %w", err)
	}
	s.sessionStore.Set(userID, encryptionKey)

	return &dusr, nil
}

func (s *service) GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(idKey).(string)
	if !ok || userID == "" {
		return "", errors.New("failed to get user id from context")
	}
	return userID, nil
}

func (s *service) SetUserIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, idKey, userID)
}
