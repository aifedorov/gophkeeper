package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	cardDomain "github.com/aifedorov/gophkeeper/internal/server/domain/secret/card"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/card/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func NewRepositoryWithTxBeginner(pool TxBeginner, querier Querier, logger *zap.Logger) interfaces.Repository {
	return &repository{
		pool:    pool,
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

func TestRepository_CreateCard(t *testing.T) {
	t.Parallel()

	t.Run("creates card successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		cardID := uuid.New()
		name := "test-card"
		encryptedNumber := []byte("encrypted-number")
		encryptedExpiredDate := []byte("encrypted-expired-date")
		encryptedCardHolderName := []byte("encrypted-card-holder-name")
		encryptedCvv := []byte("encrypted-cvv")
		encryptedNotes := []byte("encrypted-notes")

		card := interfaces.RepositoryCard{
			ID:                      cardID.String(),
			Name:                    name,
			EncryptedNumber:         encryptedNumber,
			EncryptedExpiredDate:    encryptedExpiredDate,
			EncryptedCardHolderName: encryptedCardHolderName,
			EncryptedCvv:            encryptedCvv,
			EncryptedNotes:          encryptedNotes,
		}

		now := time.Now()
		expectedDBCard := Card{
			ID:                    cardID,
			UserID:                userID,
			Name:                  name,
			EncryptedNumber:       encryptedNumber,
			EncryptedExpiredDate:  encryptedExpiredDate,
			ExpiredCardHolderName: encryptedCardHolderName,
			EncryptedCvv:          encryptedCvv,
			EncryptedNotes:        encryptedNotes,
			CreatedAt:             now,
			Version:               1,
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			CreateCard(ctx, CreateCardParams{
				ID:                    cardID,
				UserID:                userID,
				Name:                  name,
				EncryptedNumber:       encryptedNumber,
				EncryptedExpiredDate:  encryptedExpiredDate,
				ExpiredCardHolderName: encryptedCardHolderName,
				EncryptedCvv:          encryptedCvv,
				EncryptedNotes:        encryptedNotes,
			}).
			Return(expectedDBCard, nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.CreateCard(ctx, userID.String(), card)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, cardID.String(), result.ID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, encryptedNumber, result.EncryptedNumber)
		assert.Equal(t, encryptedExpiredDate, result.EncryptedExpiredDate)
		assert.Equal(t, encryptedCardHolderName, result.EncryptedCardHolderName)
		assert.Equal(t, encryptedCvv, result.EncryptedCvv)
		assert.Equal(t, encryptedNotes, result.EncryptedNotes)
	})

	t.Run("returns ErrNameExists on conflict", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		cardID := uuid.New()
		card := interfaces.RepositoryCard{
			ID:                      cardID.String(),
			Name:                    "existing-name",
			EncryptedNumber:         []byte("number"),
			EncryptedExpiredDate:    []byte("date"),
			EncryptedCardHolderName: []byte("holder"),
			EncryptedCvv:            []byte("cvv"),
			EncryptedNotes:          []byte("notes"),
		}

		pgErr := &pgconn.PgError{
			Code: pgerrcode.UniqueViolation,
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			CreateCard(ctx, gomock.Any()).
			Return(Card{}, pgErr)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.CreateCard(ctx, userID.String(), card)

		assert.ErrorIs(t, err, cardDomain.ErrNameExists)
		assert.Nil(t, result)
	})

	t.Run("returns error on database failure", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		cardID := uuid.New()
		card := interfaces.RepositoryCard{
			ID:                      cardID.String(),
			Name:                    "test",
			EncryptedNumber:         []byte("number"),
			EncryptedExpiredDate:    []byte("date"),
			EncryptedCardHolderName: []byte("holder"),
			EncryptedCvv:            []byte("cvv"),
			EncryptedNotes:          []byte("notes"),
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			CreateCard(ctx, gomock.Any()).
			Return(Card{}, sql.ErrConnDone)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.CreateCard(ctx, userID.String(), card)

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error on invalid user id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		card := interfaces.RepositoryCard{
			ID:                      uuid.New().String(),
			Name:                    "test",
			EncryptedNumber:         []byte("number"),
			EncryptedExpiredDate:    []byte("date"),
			EncryptedCardHolderName: []byte("holder"),
			EncryptedCvv:            []byte("cvv"),
			EncryptedNotes:          []byte("notes"),
		}

		mockQuerier := NewMockQuerier(ctrl)
		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.CreateCard(ctx, "invalid-uuid", card)

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error on invalid card id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		card := interfaces.RepositoryCard{
			ID:                      "invalid-uuid",
			Name:                    "test",
			EncryptedNumber:         []byte("number"),
			EncryptedExpiredDate:    []byte("date"),
			EncryptedCardHolderName: []byte("holder"),
			EncryptedCvv:            []byte("cvv"),
			EncryptedNotes:          []byte("notes"),
		}

		mockQuerier := NewMockQuerier(ctrl)
		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.CreateCard(ctx, userID.String(), card)

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRepository_ListCards(t *testing.T) {
	t.Parallel()

	t.Run("lists cards successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		card1ID := uuid.New()
		card2ID := uuid.New()

		now := time.Now()
		expectedDBCards := []Card{
			{
				ID:                    card1ID,
				UserID:                userID,
				Name:                  "card-1",
				EncryptedNumber:       []byte("number1"),
				EncryptedExpiredDate:  []byte("date1"),
				ExpiredCardHolderName: []byte("holder1"),
				EncryptedCvv:          []byte("cvv1"),
				EncryptedNotes:        []byte("notes1"),
				CreatedAt:             now,
				Version:               1,
			},
			{
				ID:                    card2ID,
				UserID:                userID,
				Name:                  "card-2",
				EncryptedNumber:       []byte("number2"),
				EncryptedExpiredDate:  []byte("date2"),
				ExpiredCardHolderName: []byte("holder2"),
				EncryptedCvv:          []byte("cvv2"),
				EncryptedNotes:        []byte("notes2"),
				CreatedAt:             now,
				Version:               1,
			},
		}

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			ListCards(ctx, userID).
			Return(expectedDBCards, nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.ListCards(ctx, userID.String())

		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, card1ID.String(), result[0].ID)
		assert.Equal(t, "card-1", result[0].Name)
		assert.Equal(t, card2ID.String(), result[1].ID)
		assert.Equal(t, "card-2", result[1].Name)
	})

	t.Run("returns empty list when no cards found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			ListCards(ctx, userID).
			Return(nil, sql.ErrNoRows)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.ListCards(ctx, userID.String())

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
			ListCards(ctx, userID).
			Return([]Card{}, nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.ListCards(ctx, userID.String())

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
			ListCards(ctx, userID).
			Return(nil, sql.ErrConnDone)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.ListCards(ctx, userID.String())

		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error on invalid user id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		result, err := repo.ListCards(ctx, "invalid-uuid")

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRepository_DeleteCard(t *testing.T) {
	t.Parallel()

	t.Run("deletes card successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		cardID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			DeleteCard(ctx, DeleteCardParams{
				ID:     cardID,
				UserID: userID,
			}).
			Return(int64(1), nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		err := repo.DeleteCard(ctx, userID.String(), cardID.String())

		require.NoError(t, err)
	})

	t.Run("returns ErrNotFound when card not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		cardID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			DeleteCard(ctx, DeleteCardParams{
				ID:     cardID,
				UserID: userID,
			}).
			Return(int64(0), nil)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		err := repo.DeleteCard(ctx, userID.String(), cardID.String())

		assert.ErrorIs(t, err, cardDomain.ErrNotFound)
	})

	t.Run("returns error on database failure", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userID := uuid.New()
		cardID := uuid.New()

		mockQuerier := NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			DeleteCard(ctx, gomock.Any()).
			Return(int64(0), sql.ErrConnDone)

		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		err := repo.DeleteCard(ctx, userID.String(), cardID.String())

		require.Error(t, err)
	})

	t.Run("returns error on invalid user id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		err := repo.DeleteCard(ctx, "invalid-uuid", uuid.New().String())

		require.Error(t, err)
	})

	t.Run("returns error on invalid card id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := NewMockQuerier(ctrl)
		repo := NewRepositoryWithQuerier(mockQuerier, logger)
		err := repo.DeleteCard(ctx, uuid.New().String(), "invalid-uuid")

		require.Error(t, err)
	})
}

func TestToInterfacesCard(t *testing.T) {
	t.Parallel()

	t.Run("converts card correctly", func(t *testing.T) {
		t.Parallel()

		cardID := uuid.New()
		userID := uuid.New()
		name := "test-card"
		encryptedNumber := []byte("encrypted-number")
		encryptedExpiredDate := []byte("encrypted-expired-date")
		encryptedCardHolderName := []byte("encrypted-card-holder-name")
		encryptedCvv := []byte("encrypted-cvv")
		encryptedNotes := []byte("encrypted-notes")

		now := time.Now()
		dbCard := Card{
			ID:                    cardID,
			UserID:                userID,
			Name:                  name,
			EncryptedNumber:       encryptedNumber,
			EncryptedExpiredDate:  encryptedExpiredDate,
			ExpiredCardHolderName: encryptedCardHolderName,
			EncryptedCvv:          encryptedCvv,
			EncryptedNotes:        encryptedNotes,
			CreatedAt:             now,
			Version:               1,
		}

		result := toInterfacesCard(dbCard)

		assert.Equal(t, cardID.String(), result.ID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, encryptedNumber, result.EncryptedNumber)
		assert.Equal(t, encryptedExpiredDate, result.EncryptedExpiredDate)
		assert.Equal(t, encryptedCardHolderName, result.EncryptedCardHolderName)
		assert.Equal(t, encryptedCvv, result.EncryptedCvv)
		assert.Equal(t, encryptedNotes, result.EncryptedNotes)
	})

	t.Run("handles empty encrypted fields", func(t *testing.T) {
		t.Parallel()

		cardID := uuid.New()
		userID := uuid.New()

		dbCard := Card{
			ID:                    cardID,
			UserID:                userID,
			Name:                  "test",
			EncryptedNumber:       []byte{},
			EncryptedExpiredDate:  []byte{},
			ExpiredCardHolderName: []byte{},
			EncryptedCvv:          []byte{},
			EncryptedNotes:        []byte{},
		}

		result := toInterfacesCard(dbCard)

		assert.Equal(t, cardID.String(), result.ID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, "test", result.Name)
		assert.Empty(t, result.EncryptedNumber)
		assert.Empty(t, result.EncryptedExpiredDate)
		assert.Empty(t, result.EncryptedCardHolderName)
		assert.Empty(t, result.EncryptedCvv)
		assert.Empty(t, result.EncryptedNotes)
	})
}
