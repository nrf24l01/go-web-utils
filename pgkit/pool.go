package pgkit

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/nrf24l01/go-web-utils/config"
)

type DB struct {
	Pool *pgxpool.Pool
	SQL  *sql.DB
}

func NewDB(ctx context.Context, pg_cfg *config.PGConfig) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(pg_cfg.GetDSN())
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(pool)

	return &DB{
		Pool: pool,
		SQL:  sqlDB,
	}, nil
}
