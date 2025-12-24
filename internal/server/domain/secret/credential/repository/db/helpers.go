package repository

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// conflictError checks if the error is a PostgreSQL unique constraint violation.
// Returns true if the error indicates a duplicate key violation (e.g., name already exists).
func conflictError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}

// notFoundError checks if the error indicates that no rows were found in the query result.
// Returns true if the error is sql.ErrNoRows.
func notFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
