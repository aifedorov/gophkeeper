package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type Repository interface {
	CreateUser(login, passHash string) (*User, error)
	GetUser(login string) (*User, error)
}

type repository struct {
	ctx     context.Context
	queries Querier
	logger  *zap.Logger
}

func NewRepository(ctx context.Context, db DBTX, logger *zap.Logger) Repository {
	return &repository{
		ctx:     ctx,
		queries: New(db),
		logger:  logger,
	}
}

func NewRepositoryWithQuerier(ctx context.Context, querier Querier, logger *zap.Logger) Repository {
	return &repository{
		ctx:     ctx,
		queries: querier,
		logger:  logger,
	}
}

func (s *repository) CreateUser(login, passHash string) (*User, error) {
	user, err := s.queries.CreateUser(s.ctx, CreateUserParams{
		Login:        login,
		PasswordHash: passHash,
	})
	if IsConflictError(err) {
		return nil, ErrLoginExists
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to create auth.proto: %w", err)
	}
	return &user, err
}

func (s *repository) GetUser(login string) (*User, error) {
	user, err := s.queries.GetUser(s.ctx, login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to get auth.proto: %w", err)
	}
	return &user, nil
}
