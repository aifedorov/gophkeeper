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

type service struct {
	ctx     context.Context
	queries *Queries
	logger  *zap.Logger
}

func NewRepository(ctx context.Context, db DBTX, logger *zap.Logger) Repository {
	return &service{
		ctx:     ctx,
		queries: New(db),
		logger:  logger,
	}
}

func (s *service) CreateUser(login, passHash string) (*User, error) {
	user, err := s.queries.CreateUser(s.ctx, CreateUserParams{
		Login:        login,
		PasswordHash: passHash,
	})
	if IsConflictError(err) {
		return nil, ErrLoginExists
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to create user: %w", err)
	}
	return &user, err
}

func (s *service) GetUser(login, passHash string) (*User, error) {
	user, err := s.queries.GetUser(s.ctx, GetUserParams{
		Login:        login,
		PasswordHash: passHash,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo: failed to get user: %w", err)
	}
	return &user, nil
}
