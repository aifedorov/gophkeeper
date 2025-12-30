package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultDBTimeout = 3 * time.Second

type Connection struct {
	dbPool *pgxpool.Pool
}

// NewConnection creates and opens a new PostgreSQL connection pool.
func NewConnection(ctx context.Context, dsn string) (*Connection, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDBTimeout)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &Connection{dbPool: dbpool}, nil
}

func (p *Connection) Close() {
	if p.dbPool != nil {
		p.dbPool.Close()
	}
}

func (p *Connection) DBPool() *pgxpool.Pool {
	return p.dbPool
}
