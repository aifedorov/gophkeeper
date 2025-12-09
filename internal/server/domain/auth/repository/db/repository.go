package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type Repository interface {
	CreateUser(ctx context.Context, login, passHash string) (*User, error)
	GetUser(ctx context.Context, login string) (*User, error)
}

type repository struct {
	queries Querier
	logger  *zap.Logger
}

func NewRepository(db DBTX, logger *zap.Logger) Repository {
	return &repository{
		queries: New(db),
		logger:  logger,
	}
}

func NewRepositoryWithQuerier(querier Querier, logger *zap.Logger) Repository {
	return &repository{
		queries: querier,
		logger:  logger,
	}
}

func (s *repository) CreateUser(ctx context.Context, login, passHash string) (*User, error) {
	user, err := s.queries.CreateUser(ctx, CreateUserParams{
		Login:        login,
		PasswordHash: passHash,
	})
	if conflictError(err) {
		return nil, ErrLoginExists
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to create auth.proto: %w", err)
	}
	return &user, err
}

func (s *repository) GetUser(ctx context.Context, login string) (*User, error) {
	user, err := s.queries.GetUser(ctx, login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to get auth.proto: %w", err)
	}
	return &user, nil
}
