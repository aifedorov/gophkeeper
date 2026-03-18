package repository

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func conflictError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}

// notFoundError checks if the error indicates that no rows were found in the query result.
// Returns true if the error is sql.ErrNoRows.
func notFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
