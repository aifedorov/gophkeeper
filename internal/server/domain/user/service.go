package user

import (
	"errors"
	"fmt"

	repository2 "github.com/aifedorov/gophkeeper/internal/server/domain/user/repository/db"
	"go.uber.org/zap"
)

type Service interface {
	Register(login, passHash string) (*User, error)
	Login(login, passHash string) (*User, error)
}

type service struct {
	repo   repository2.Repository
	logger *zap.Logger
}

func NewService(repo repository2.Repository, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Register(login, passHash string) (*User, error) {
	user, err := s.repo.CreateUser(login, passHash)
	if errors.Is(err, repository2.ErrLoginExists) {
		return nil, ErrLoginExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return toDomainUser(user), nil
}

func (s *service) Login(login, passHash string) (*User, error) {
	user, err := s.repo.GetUser(login, passHash)
	if errors.Is(err, repository2.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return toDomainUser(user), nil
}
