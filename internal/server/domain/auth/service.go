package auth

import (
	"context"
	"errors"
	"fmt"

	repository "github.com/aifedorov/gophkeeper/internal/server/domain/auth/repository/db"
	"github.com/aifedorov/gophkeeper/pkg/validator"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the interface for auth domain operations.
type Service interface {
	// Register creates a new auth account with the provided login and password.
	// The password is hashed before storage. Returns ErrLoginExists if the login is already taken.
	Register(ctx context.Context, login, password string) (*User, error)
	// Login authenticates a auth with the provided credentials.
	// Returns ErrUserNotFound if the auth doesn't exist or if the password is incorrect.
	Login(ctx context.Context, login, password string) (*User, error)
}

type service struct {
	repo   repository.Repository
	logger *zap.Logger
}

// NewService creates a new instance of the auth service with the provided repository and logger.
// It initializes the service that handles auth registration and authentication business logic.
func NewService(repo repository.Repository, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Register(ctx context.Context, login, password string) (*User, error) {
	if err := validator.ValidateLogin(login); err != nil {
		return nil, fmt.Errorf("invalid login: %w", err)
	}
	if err := validator.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, login, string(hashedPassword))
	if errors.Is(err, repository.ErrLoginExists) {
		return nil, ErrLoginExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create auth.proto: %w", err)
	}
	return toDomainUser(user), nil
}

func (s *service) Login(ctx context.Context, login, password string) (*User, error) {
	if err := validator.ValidateLogin(login); err != nil {
		return nil, fmt.Errorf("invalid login: %w", err)
	}
	if err := validator.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	user, err := s.repo.GetUser(ctx, login)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get auth.proto: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("failed to compare hash and password: %w", err)
	}

	return toDomainUser(user), nil
}
