// Package repository provides database operations for card management.
package repository

import (
	"context"
	"fmt"

	cardDomain "github.com/aifedorov/gophkeeper/internal/server/domain/secret/card"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/card/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// repository implements the interfaces.Repository for card persistence.
type repository struct {
	pool    TxBeginner
	queries Querier
	logger  *zap.Logger
}

// NewRepository creates a new card repository with database connection and logger.
// The pool is used for transaction management while db is used for query execution.
func NewRepository(pool *pgxpool.Pool, db DBTX, logger *zap.Logger) interfaces.Repository {
	return &repository{
		pool:    pool,
		queries: newQuerier(db),
		logger:  logger,
	}
}

func (r *repository) CreateCard(ctx context.Context, userID string, card interfaces.RepositoryCard) (*interfaces.RepositoryCard, error) {
	r.logger.Debug("repo: creating card",
		zap.String("user_id", userID),
		zap.String("name", card.Name))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	id, err := uuid.Parse(card.ID)
	if err != nil {
		r.logger.Error("repo: failed to parse card id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse card id: %w", err)
	}

	dbCard, err := r.queries.CreateCard(ctx, CreateCardParams{
		ID:                    id,
		UserID:                userUUID,
		Name:                  card.Name,
		EncryptedNumber:       card.EncryptedNumber,
		EncryptedExpiredDate:  card.EncryptedExpiredDate,
		ExpiredCardHolderName: card.EncryptedCardHolderName,
		EncryptedCvv:          card.EncryptedCvv,
		EncryptedNotes:        card.EncryptedNotes,
	})
	if conflictError(err) {
		r.logger.Debug("repo: card name already exists", zap.String("name", card.Name))
		return nil, cardDomain.ErrNameExists
	}
	if err != nil {
		r.logger.Error("repo: failed to create card", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to create card: %w", err)
	}
	result := toInterfacesCard(dbCard)
	r.logger.Debug("repo: card created successfully", zap.String("id", result.ID))
	return &result, nil
}

func (r *repository) ListCards(ctx context.Context, userID string) ([]interfaces.RepositoryCard, error) {
	r.logger.Debug("repo: listing cards", zap.String("user_id", userID))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	dbCards, err := r.queries.ListCards(ctx, userUUID)
	if notFoundError(err) {
		r.logger.Debug("repo: no cards found", zap.String("user_id", userID))
		return []interfaces.RepositoryCard{}, nil
	}
	if err != nil {
		r.logger.Error("repo: failed to list cards", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to list cards: %w", err)
	}
	result := make([]interfaces.RepositoryCard, len(dbCards))
	for i, card := range dbCards {
		result[i] = toInterfacesCard(card)
	}
	r.logger.Debug("repo: cards listed successfully", zap.Int("count", len(result)))
	return result, nil
}

func (r *repository) UpdateCard(ctx context.Context, userID string, card interfaces.RepositoryCard) (*interfaces.RepositoryCard, error) {
	r.logger.Debug("repo: updating card",
		zap.String("user_id", userID),
		zap.String("id", card.ID))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	id, err := uuid.Parse(card.ID)
	if err != nil {
		r.logger.Error("repo: failed to parse card id", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to parse card id: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.Error("repo: failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to begin transaction: %w", err)
	}
	defer func() {
		r.logger.Debug("repo: rollback transaction for update")
		_ = tx.Rollback(ctx)
	}()

	txQuery := r.queries.WithTx(tx)
	cardRepo, err := txQuery.GetCardForUpdate(ctx, GetCardForUpdateParams{
		ID:     id,
		UserID: userUUID,
	})
	if notFoundError(err) {
		r.logger.Debug("repo: card not found for update", zap.String("id", card.ID))
		return nil, cardDomain.ErrNotFound
	}
	if err != nil {
		r.logger.Error("repo: failed to get card for update", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to get card for update: %w", err)
	}

	if cardRepo.Version != card.Version {
		r.logger.Debug("repo: version conflict",
			zap.String("id", card.ID),
			zap.Int64("db_version", cardRepo.Version),
			zap.Int64("client_version", card.Version),
		)
		return nil, cardDomain.ErrVersionConflict
	}

	dbCard, err := txQuery.UpdateCard(ctx, UpdateCardParams{
		ID:                    id,
		UserID:                userUUID,
		Version:               card.Version,
		Name:                  card.Name,
		EncryptedNumber:       card.EncryptedNumber,
		EncryptedExpiredDate:  card.EncryptedExpiredDate,
		ExpiredCardHolderName: card.EncryptedCardHolderName,
		EncryptedCvv:          card.EncryptedCvv,
		EncryptedNotes:        card.EncryptedNotes,
	})
	if notFoundError(err) {
		r.logger.Debug("repo: version conflict", zap.String("id", card.ID))
		return nil, cardDomain.ErrVersionConflict
	}
	if conflictError(err) {
		r.logger.Debug("repo: card name already exists", zap.String("name", card.Name))
		return nil, cardDomain.ErrNameExists
	}
	if err != nil {
		r.logger.Error("repo: failed to update card", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to update card: %w", err)
	}
	result := toInterfacesCard(dbCard)
	r.logger.Debug("repo: card updated successfully", zap.String("id", result.ID))

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("repo: failed to commit transaction", zap.Error(err))
		return nil, fmt.Errorf("repo: failed to commit transaction: %w", err)
	}

	r.logger.Debug("repo: transaction committed successfully")
	return &result, nil
}

func (r *repository) DeleteCard(ctx context.Context, userID, id string) error {
	r.logger.Debug("repo: deleting card",
		zap.String("user_id", userID),
		zap.String("id", id))

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("repo: failed to parse user id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse user id: %w", err)
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("repo: failed to parse card id", zap.Error(err))
		return fmt.Errorf("repo: failed to parse card id: %w", err)
	}

	rows, err := r.queries.DeleteCard(ctx, DeleteCardParams{
		ID:     idUUID,
		UserID: userUUID,
	})
	if rows == 0 {
		r.logger.Debug("repo: card not found for deletion", zap.String("id", id))
		return cardDomain.ErrNotFound
	}
	if err != nil {
		r.logger.Error("repo: failed to delete card", zap.Error(err))
		return fmt.Errorf("repo: failed to delete card: %w", err)
	}
	r.logger.Debug("repo: card deleted successfully", zap.String("id", id))
	return nil
}

func toInterfacesCard(card Card) interfaces.RepositoryCard {
	return interfaces.RepositoryCard{
		ID:                      card.ID.String(),
		UserID:                  card.UserID.String(),
		Name:                    card.Name,
		EncryptedNumber:         card.EncryptedNumber,
		EncryptedExpiredDate:    card.EncryptedExpiredDate,
		EncryptedCardHolderName: card.ExpiredCardHolderName,
		EncryptedCvv:            card.EncryptedCvv,
		EncryptedNotes:          card.EncryptedNotes,
		Version:                 card.Version,
	}
}
