package user

import (
	"errors"
	"fmt"

	repository "github.com/aifedorov/gophkeeper/internal/server/domain/user/repository/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the interface for user domain operations.
type Service interface {
	// Register creates a new user account with the provided login and password.
	// The password is hashed before storage. Returns ErrLoginExists if the login is already taken.
	Register(login, password string) (*User, error)
	// Login authenticates a user with the provided credentials.
	// Returns ErrUserNotFound if the user doesn't exist or if the password is incorrect.
	Login(login, password string) (*User, error)
}

type service struct {
	repo   repository.Repository
	logger *zap.Logger
}

// NewService creates a new instance of the user service with the provided repository and logger.
// It initializes the service that handles user registration and authentication business logic.
func NewService(repo repository.Repository, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Register(login, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.CreateUser(login, string(hashedPassword))
	if errors.Is(err, repository.ErrLoginExists) {
		return nil, ErrLoginExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create auth.proto: %w", err)
	}
	return toDomainUser(user), nil
}

func (s *service) Login(login, password string) (*User, error) {
	user, err := s.repo.GetUser(login)
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
