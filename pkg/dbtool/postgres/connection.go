package postgres

// This package is created for providing helper functions when working with database

// Currently I've only wrote for connection initializing and transaction wrapper.

// The reason both package "pgx" and package "sqlx" be used is each one has its own strength.
// "pgx": better connection management, help us leverage Postgres specific features like LISTEN/NOTIFY, COPY, etc...
// "sqlx": better for scan struct and parameters binding, also I've worked with "sqlx" more than "pgx"

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func GetPgxConn(dsn string, cfg SingleConnConfig) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.InitTimeout)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func GetPgxPool(dsn string, cfg PoolConfig) (*pgxpool.Pool, error) {
	pgxPoolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	if cfg.MaxConnIdleTime != 0 {
		pgxPoolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	}
	if cfg.MinConns != 0 {
		pgxPoolCfg.MinConns = int32(cfg.MinConns)
	}
	if cfg.MaxConns != 0 {
		pgxPoolCfg.MaxConns = int32(cfg.MaxConns)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.InitTimeout)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, pgxPoolCfg)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// *sqlx.DB can be considered as equivalent to *pgxpool.Pool since both manage connections
// to database
func GetSqlxDB(dsn string, cfg PoolConfig) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.InitTimeout)
	defer cancel()
	pool, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
