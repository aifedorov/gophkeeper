package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/domain/auth/interfaces"
	"go.uber.org/zap"
)

type repository struct {
	queries Querier
	logger  *zap.Logger
}

func NewRepository(db DBTX, logger *zap.Logger) interfaces.Repository {
	return &repository{
		queries: New(db),
		logger:  logger,
	}
}

func NewRepositoryWithQuerier(querier Querier, logger *zap.Logger) interfaces.Repository {
	return &repository{
		queries: querier,
		logger:  logger,
	}
}

func (s *repository) CreateUser(ctx context.Context, user interfaces.RepositoryUser, passwordHash string) (interfaces.RepositoryUser, error) {
	dbUser, err := s.queries.CreateUser(ctx, CreateUserParams{
		Login:        user.Login,
		PasswordHash: passwordHash,
		Salt:         user.Salt,
	})
	if conflictError(err) {
		return interfaces.RepositoryUser{}, auth.ErrLoginExists
	}
	if err != nil {
		return interfaces.RepositoryUser{}, fmt.Errorf("repo: failed to create user: %w", err)
	}
	return toInterfacesUser(dbUser), err
}

func (s *repository) GetUser(ctx context.Context, login string) (interfaces.RepositoryUser, error) {
	dbUser, err := s.queries.GetUser(ctx, login)
	if errors.Is(err, sql.ErrNoRows) {
		return interfaces.RepositoryUser{}, auth.ErrUserNotFound
	}
	if err != nil {
		return interfaces.RepositoryUser{}, fmt.Errorf("repo: failed to get user: %w", err)
	}
	return toInterfacesUser(dbUser), nil
}

func toInterfacesUser(user User) interfaces.RepositoryUser {
	return interfaces.RepositoryUser{
		ID:           user.ID.String(),
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		Salt:         user.Salt,
	}
}
