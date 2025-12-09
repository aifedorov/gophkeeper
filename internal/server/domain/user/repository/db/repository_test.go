package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestNewRepository(t *testing.T) {
	t.Parallel()

	t.Run("creates repository", func(t *testing.T) {
		t.Parallel()

		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := NewRepositoryWithQuerier(mockQuerier, logger)

		require.NotNil(t, repo)
	})
}

func TestRepository_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("creates user successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		login := "testuser"
		passwordHash := "hashedpassword"

		expectedUser := User{
			ID:           userID,
			Login:        login,
			PasswordHash: passwordHash,
			CreatedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			CreateUser(ctx, CreateUserParams{
				Login:        login,
				PasswordHash: passwordHash,
			}).
			Return(expectedUser, nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		user, err := repo.CreateUser(ctx, login, passwordHash)

		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, login, user.Login)
		assert.Equal(t, passwordHash, user.PasswordHash)
	})

	t.Run("returns ErrLoginExists on conflict", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pgErr := &pgconn.PgError{
			Code: pgerrcode.UniqueViolation,
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			CreateUser(ctx, CreateUserParams{
				Login:        "existinguser",
				PasswordHash: "password",
			}).
			Return(User{}, pgErr)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		user, err := repo.CreateUser(ctx, "existinguser", "password")

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrLoginExists)
		assert.Nil(t, user)
	})

	t.Run("returns error on database error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbErr := errors.New("database connection failed")

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			CreateUser(ctx, CreateUserParams{
				Login:        "testuser",
				PasswordHash: "password",
			}).
			Return(User{}, dbErr)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		user, err := repo.CreateUser(ctx, "testuser", "password")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create auth.proto")
		assert.Nil(t, user)
	})
}

func TestRepository_GetUser(t *testing.T) {
	t.Parallel()

	t.Run("gets user successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		login := "testuser"
		passwordHash := "hashedpassword"

		expectedUser := User{
			ID:           userID,
			Login:        login,
			PasswordHash: passwordHash,
			CreatedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			GetUser(ctx, login).
			Return(expectedUser, nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		user, err := repo.GetUser(ctx, login)

		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, login, user.Login)
		assert.Equal(t, passwordHash, user.PasswordHash)
	})

	t.Run("returns ErrUserNotFound when user not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			GetUser(ctx, "nonexistent").
			Return(User{}, sql.ErrNoRows)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		user, err := repo.GetUser(ctx, "nonexistent")

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, user)
	})

	t.Run("returns error on database error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbErr := errors.New("database connection failed")

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			GetUser(ctx, "testuser").
			Return(User{}, dbErr)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		user, err := repo.GetUser(ctx, "testuser")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get auth.proto")
		assert.Nil(t, user)
	})
}
