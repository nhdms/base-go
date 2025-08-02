package postgres

import "time"

// Config for single connection
type SingleConnConfig struct {
	InitTimeout time.Duration
}

// Config for connection pool
type PoolConfig struct {
	InitTimeout     time.Duration
	MaxConnIdleTime time.Duration
	MinConns        int
	MaxConns        int
}
