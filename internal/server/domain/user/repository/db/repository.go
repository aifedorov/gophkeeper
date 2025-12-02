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
	GetUser(login, passHash string) (*User, error)
}

type repository struct {
	ctx     context.Context
	queries *Queries
	logger  *zap.Logger
}

func NewRepository(ctx context.Context, db DBTX, logger *zap.Logger) Repository {
	return &repository{
		ctx:     ctx,
		queries: New(db),
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

func (s *repository) GetUser(login, passHash string) (*User, error) {
	user, err := s.queries.GetUser(s.ctx, GetUserParams{
		Login:        login,
		PasswordHash: passHash,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to get auth.proto: %w", err)
	}
	return &user, nil
}
