package repository

import (
	"github.com/jackc/pgx/v5"
)

// queriesWrapper wraps the sqlc-generated Queries to implement the Querier interface
type queriesWrapper struct {
	*Queries
}

// WithTx returns a new Querier that uses the given transaction
func (q *queriesWrapper) WithTx(tx pgx.Tx) Querier {
	return &queriesWrapper{Queries: q.Queries.WithTx(tx)}
}

// Ensure queriesWrapper implements Querier
var _ Querier = (*queriesWrapper)(nil)

// newQuerier creates a new Querier from a DBTX
func newQuerier(db DBTX) Querier {
	return &queriesWrapper{Queries: New(db)}
}
