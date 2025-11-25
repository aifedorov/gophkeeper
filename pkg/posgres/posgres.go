package posgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultDBTimeout = 3 * time.Second

type PostgresConnection struct {
	dbPool *pgxpool.Pool
	ctx    context.Context
	dsn    string
}

func NewPosgresConnection(ctx context.Context, dsn string) *PostgresConnection {
	return &PostgresConnection{
		ctx: ctx,
		dsn: dsn,
	}
}

func (p *PostgresConnection) Open() error {
	ctx, cancel := context.WithTimeout(p.ctx, defaultDBTimeout)
	defer cancel()

	dbpool, err := pgxpool.New(context.Background(), p.dsn)
	if err != nil {
		return err
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	p.dbPool = dbpool
	return nil
}

func (p *PostgresConnection) Close() {
	p.dbPool.Close()
}

func (p *PostgresConnection) DBPool() *pgxpool.Pool {
	return p.dbPool
}
