package auth

import (
	"context"
	"testing"

	"github.com/aifedorov/gophkeeper/internal/client/infrastructure/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewRepository(t *testing.T) {
	t.Parallel()

	t.Run("creates repository", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		storage := memory.NewStorage()

		repo := NewRepository(ctx, logger, storage)

		require.NotNil(t, repo)
	})
}

func TestRepository_Save(t *testing.T) {
	t.Parallel()

	t.Run("saves session successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		storage := memory.NewStorage()
		repo := NewRepository(ctx, logger, storage)

		session := Session{
			User: User{
				ID:    "auth-id-123",
				Login: "testuser",
			},
			AccessToken: "token-xyz",
		}

		err := repo.Save(session)
		require.NoError(t, err)

		// Verify session was saved
		loaded, err := repo.Load()
		require.NoError(t, err)
		assert.Equal(t, session, loaded)
	})
}

func TestRepository_Load(t *testing.T) {
	t.Parallel()

	t.Run("loads session successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		storage := memory.NewStorage()
		repo := NewRepository(ctx, logger, storage)

		expectedSession := Session{
			User: User{
				ID:    "auth-id-123",
				Login: "testuser",
			},
			AccessToken: "token-xyz",
		}

		err := repo.Save(expectedSession)
		require.NoError(t, err)

		session, err := repo.Load()
		require.NoError(t, err)
		assert.Equal(t, expectedSession, session)
	})

	t.Run("returns error when session not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		storage := memory.NewStorage()
		repo := NewRepository(ctx, logger, storage)

		_, err := repo.Load()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrSessionNotFound)
	})
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()

	t.Run("deletes session successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		storage := memory.NewStorage()
		repo := NewRepository(ctx, logger, storage)

		session := Session{
			User: User{
				ID:    "auth-id-123",
				Login: "testuser",
			},
			AccessToken: "token-xyz",
		}

		err := repo.Save(session)
		require.NoError(t, err)

		err = repo.Delete()
		require.NoError(t, err)

		// Verify session was deleted
		_, err = repo.Load()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrSessionNotFound)
	})
}

func TestToDomainUser(t *testing.T) {
	t.Parallel()

	t.Run("converts memory auth to domain auth", func(t *testing.T) {
		t.Parallel()

		memoryUser := memory.User{
			ID:    "auth-id-123",
			Login: "testuser",
		}

		domainUser := toDomainUser(memoryUser)

		assert.Equal(t, "auth-id-123", domainUser.ID)
		assert.Equal(t, "testuser", domainUser.Login)
	})
}

func TestToDomainSession(t *testing.T) {
	t.Parallel()

	t.Run("converts memory session to domain session", func(t *testing.T) {
		t.Parallel()

		memorySession := memory.Session{
			User: memory.User{
				ID:    "auth-id-123",
				Login: "testuser",
			},
			AccessToken: "token-xyz",
		}

		domainSession := toDomainSession(memorySession)

		assert.Equal(t, "auth-id-123", domainSession.User.ID)
		assert.Equal(t, "testuser", domainSession.User.Login)
		assert.Equal(t, "token-xyz", domainSession.AccessToken)
	})
}

func TestToMemorySession(t *testing.T) {
	t.Parallel()

	t.Run("converts domain session to memory session", func(t *testing.T) {
		t.Parallel()

		domainSession := Session{
			User: User{
				ID:    "auth-id-123",
				Login: "testuser",
			},
			AccessToken: "token-xyz",
		}

		memorySession := toMemorySession(domainSession)

		assert.Equal(t, "auth-id-123", memorySession.User.ID)
		assert.Equal(t, "testuser", memorySession.User.Login)
		assert.Equal(t, "token-xyz", memorySession.AccessToken)
	})
}

func TestToMemoryUser(t *testing.T) {
	t.Parallel()

	t.Run("converts domain auth to memory auth", func(t *testing.T) {
		t.Parallel()

		domainUser := User{
			ID:    "auth-id-123",
			Login: "testuser",
		}

		memoryUser := toMemoryUser(domainUser)

		assert.Equal(t, "auth-id-123", memoryUser.ID)
		assert.Equal(t, "testuser", memoryUser.Login)
	})
}
