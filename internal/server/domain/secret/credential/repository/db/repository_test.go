package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	credentialDomain "github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/credential/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func NewRepositoryWithQuerier(querier Querier, logger *zap.Logger) interfaces.Repository {
	return &repository{
		queries: querier,
		logger:  logger,
	}
}

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

func TestRepository_CreateCredential(t *testing.T) {
	t.Parallel()

	t.Run("creates credential successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		credID := uuid.New()
		name := "test-credential"
		encryptedLogin := []byte("encrypted-login")
		encryptedPassword := []byte("encrypted-password")
		encryptedNotes := []byte("encrypted-notes")

		credential := interfaces.RepositoryCredential{
			Name:              name,
			EncryptedLogin:    encryptedLogin,
			EncryptedPassword: encryptedPassword,
			EncryptedNotes:    encryptedNotes,
		}

		expectedDBCred := Credential{
			ID:                credID,
			UserID:            userID,
			Name:              name,
			Encryptedlogin:    encryptedLogin,
			Encryptedpassword: encryptedPassword,
			Encryptednotes:    encryptedNotes,
			CreatedAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			CreateCredential(ctx, CreateCredentialParams{
				UserID:            userID,
				Name:              name,
				Encryptedlogin:    encryptedLogin,
				Encryptedpassword: encryptedPassword,
				Encryptednotes:    encryptedNotes,
			}).
			Return(expectedDBCred, nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.CreateCredential(ctx, userID.String(), credential)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, credID.String(), result.ID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, encryptedLogin, result.EncryptedLogin)
		assert.Equal(t, encryptedPassword, result.EncryptedPassword)
		assert.Equal(t, encryptedNotes, result.EncryptedNotes)
	})

	t.Run("returns ErrNameExists on conflict", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		credential := interfaces.RepositoryCredential{
			Name:              "existing-name",
			EncryptedLogin:    []byte("login"),
			EncryptedPassword: []byte("password"),
			EncryptedNotes:    []byte("notes"),
		}

		pgErr := &pgconn.PgError{
			Code: pgerrcode.UniqueViolation,
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			CreateCredential(ctx, gomock.Any()).
			Return(Credential{}, pgErr)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.CreateCredential(ctx, userID.String(), credential)

		assert.ErrorIs(t, err, credentialDomain.ErrNameExists)
		assert.Nil(t, result)
	})

	t.Run("returns error on database failure", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		credential := interfaces.RepositoryCredential{
			Name:              "test",
			EncryptedLogin:    []byte("login"),
			EncryptedPassword: []byte("password"),
			EncryptedNotes:    []byte("notes"),
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			CreateCredential(ctx, gomock.Any()).
			Return(Credential{}, sql.ErrConnDone)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.CreateCredential(ctx, userID.String(), credential)

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRepository_ListCredentials(t *testing.T) {
	t.Parallel()

	t.Run("lists credentials successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		cred1ID := uuid.New()
		cred2ID := uuid.New()

		expectedDBCreds := []Credential{
			{
				ID:                cred1ID,
				UserID:            userID,
				Name:              "credential-1",
				Encryptedlogin:    []byte("login1"),
				Encryptedpassword: []byte("password1"),
				Encryptednotes:    []byte("notes1"),
				CreatedAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
			{
				ID:                cred2ID,
				UserID:            userID,
				Name:              "credential-2",
				Encryptedlogin:    []byte("login2"),
				Encryptedpassword: []byte("password2"),
				Encryptednotes:    []byte("notes2"),
				CreatedAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			ListCredentials(ctx, userID).
			Return(expectedDBCreds, nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.ListCredentials(ctx, userID.String())

		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, cred1ID.String(), result[0].ID)
		assert.Equal(t, "credential-1", result[0].Name)
		assert.Equal(t, cred2ID.String(), result[1].ID)
		assert.Equal(t, "credential-2", result[1].Name)
	})

	t.Run("returns empty list when no credentials found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			ListCredentials(ctx, userID).
			Return(nil, sql.ErrNoRows)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.ListCredentials(ctx, userID.String())

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("returns empty list when database returns empty slice", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			ListCredentials(ctx, userID).
			Return([]Credential{}, nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.ListCredentials(ctx, userID.String())

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("returns error on database failure", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			ListCredentials(ctx, userID).
			Return(nil, sql.ErrConnDone)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.ListCredentials(ctx, userID.String())

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRepository_UpdateCredential(t *testing.T) {
	t.Parallel()

	t.Run("updates credential successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		credID := uuid.New()
		name := "updated-credential"
		encryptedLogin := []byte("updated-login")
		encryptedPassword := []byte("updated-password")
		encryptedNotes := []byte("updated-notes")

		credential := interfaces.RepositoryCredential{
			ID:                credID.String(),
			Name:              name,
			EncryptedLogin:    encryptedLogin,
			EncryptedPassword: encryptedPassword,
			EncryptedNotes:    encryptedNotes,
		}

		expectedDBCred := Credential{
			ID:                credID,
			UserID:            userID,
			Name:              name,
			Encryptedlogin:    encryptedLogin,
			Encryptedpassword: encryptedPassword,
			Encryptednotes:    encryptedNotes,
			UpdatedAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			UpdateCredential(ctx, UpdateCredentialParams{
				ID:                credID,
				UserID:            userID,
				Name:              name,
				Encryptedlogin:    encryptedLogin,
				Encryptedpassword: encryptedPassword,
				Encryptednotes:    encryptedNotes,
			}).
			Return(expectedDBCred, nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.UpdateCredential(ctx, userID.String(), credential)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, credID.String(), result.ID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, encryptedLogin, result.EncryptedLogin)
		assert.Equal(t, encryptedPassword, result.EncryptedPassword)
		assert.Equal(t, encryptedNotes, result.EncryptedNotes)
	})

	t.Run("returns ErrNotFound when credential not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		credID := uuid.New()

		credential := interfaces.RepositoryCredential{
			ID:                credID.String(),
			Name:              "test",
			EncryptedLogin:    []byte("login"),
			EncryptedPassword: []byte("password"),
			EncryptedNotes:    []byte("notes"),
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			UpdateCredential(ctx, gomock.Any()).
			Return(Credential{}, sql.ErrNoRows)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.UpdateCredential(ctx, userID.String(), credential)

		assert.ErrorIs(t, err, credentialDomain.ErrNotFound)
		assert.Nil(t, result)
	})

	t.Run("returns error on invalid credential ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()

		credential := interfaces.RepositoryCredential{
			ID:                "invalid-uuid",
			Name:              "test",
			EncryptedLogin:    []byte("login"),
			EncryptedPassword: []byte("password"),
			EncryptedNotes:    []byte("notes"),
		}

		mockQuerier := NewMockQuerier(ctrl)
		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.UpdateCredential(ctx, userID.String(), credential)

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error on database failure", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		credID := uuid.New()

		credential := interfaces.RepositoryCredential{
			ID:                credID.String(),
			Name:              "test",
			EncryptedLogin:    []byte("login"),
			EncryptedPassword: []byte("password"),
			EncryptedNotes:    []byte("notes"),
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			UpdateCredential(ctx, gomock.Any()).
			Return(Credential{}, sql.ErrConnDone)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.UpdateCredential(ctx, userID.String(), credential)

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRepository_DeleteCredential(t *testing.T) {
	t.Parallel()

	t.Run("deletes credential successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		credID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			DeleteCredential(ctx, DeleteCredentialParams{
				ID:     credID,
				UserID: userID,
			}).
			Return(int64(1), nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		err := repo.DeleteCredential(ctx, userID.String(), credID.String())

		require.NoError(t, err)
	})

	t.Run("returns ErrNotFound when credential not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		credID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			DeleteCredential(ctx, DeleteCredentialParams{
				ID:     credID,
				UserID: userID,
			}).
			Return(int64(0), nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		err := repo.DeleteCredential(ctx, userID.String(), credID.String())

		assert.ErrorIs(t, err, credentialDomain.ErrNotFound)
	})

	t.Run("returns error on database failure", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		credID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			DeleteCredential(ctx, gomock.Any()).
			Return(int64(0), sql.ErrConnDone)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		err := repo.DeleteCredential(ctx, userID.String(), credID.String())

		require.Error(t, err)
	})
}

func TestToInterfacesCredential(t *testing.T) {
	t.Parallel()

	t.Run("converts credential correctly", func(t *testing.T) {
		t.Parallel()

		credID := uuid.New()
		userID := uuid.New()
		name := "test-credential"
		encryptedLogin := []byte("encrypted-login")
		encryptedPassword := []byte("encrypted-password")
		encryptedNotes := []byte("encrypted-notes")

		dbCred := Credential{
			ID:                credID,
			UserID:            userID,
			Name:              name,
			Encryptedlogin:    encryptedLogin,
			Encryptedpassword: encryptedPassword,
			Encryptednotes:    encryptedNotes,
			CreatedAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		result := toInterfacesCredential(dbCred)

		assert.Equal(t, credID.String(), result.ID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, encryptedLogin, result.EncryptedLogin)
		assert.Equal(t, encryptedPassword, result.EncryptedPassword)
		assert.Equal(t, encryptedNotes, result.EncryptedNotes)
	})

	t.Run("handles empty encrypted fields", func(t *testing.T) {
		t.Parallel()

		credID := uuid.New()
		userID := uuid.New()

		dbCred := Credential{
			ID:                credID,
			UserID:            userID,
			Name:              "test",
			Encryptedlogin:    []byte{},
			Encryptedpassword: []byte{},
			Encryptednotes:    []byte{},
		}

		result := toInterfacesCredential(dbCred)

		assert.Equal(t, credID.String(), result.ID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, "test", result.Name)
		assert.Empty(t, result.EncryptedLogin)
		assert.Empty(t, result.EncryptedPassword)
		assert.Empty(t, result.EncryptedNotes)
	})
}
