package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=tx.go -destination=mock_tx_test.go -package=repository
//go:generate mockgen -package=repository -destination=mock_pgx_tx_test.go github.com/jackc/pgx/v5 Tx

// TxBeginner is an interface for beginning transactions.
// Implemented by *pgxpool.Pool.
type TxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}
