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
	s.logger.Debug("repo: creating user", zap.String("login", user.Login))

	dbUser, err := s.queries.CreateUser(ctx, CreateUserParams{
		Login:        user.Login,
		PasswordHash: passwordHash,
		Salt:         []byte(user.Salt),
	})
	if conflictError(err) {
		s.logger.Debug("repo: user already exists", zap.String("login", user.Login))
		return interfaces.RepositoryUser{}, auth.ErrLoginExists
	}
	if err != nil {
		s.logger.Error("repo: failed to create user", zap.Error(err))
		return interfaces.RepositoryUser{}, fmt.Errorf("repo: failed to create user: %w", err)
	}

	s.logger.Debug("repo: user created successfully", zap.String("user_id", dbUser.ID.String()))
	return toInterfacesUser(dbUser), err
}

func (s *repository) GetUser(ctx context.Context, login string) (interfaces.RepositoryUser, error) {
	s.logger.Debug("repo: getting user", zap.String("login", login))

	dbUser, err := s.queries.GetUser(ctx, login)
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.Debug("repo: user not found", zap.String("login", login))
		return interfaces.RepositoryUser{}, auth.ErrUserNotFound
	}
	if err != nil {
		s.logger.Error("repo: failed to get user", zap.Error(err))
		return interfaces.RepositoryUser{}, fmt.Errorf("repo: failed to get user: %w", err)
	}

	s.logger.Debug("repo: user found", zap.String("user_id", dbUser.ID.String()))
	return toInterfacesUser(dbUser), nil
}

func toInterfacesUser(user User) interfaces.RepositoryUser {
	return interfaces.RepositoryUser{
		ID:           user.ID.String(),
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		Salt:         string(user.Salt),
	}
}
