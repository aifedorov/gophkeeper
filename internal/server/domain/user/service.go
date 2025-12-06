package user

import (
	"errors"
	"fmt"

	repository "github.com/aifedorov/gophkeeper/internal/server/domain/user/repository/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(login, password string) (*User, error)
	Login(login, password string) (*User, error)
}

type service struct {
	repo   repository.Repository
	logger *zap.Logger
}

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
