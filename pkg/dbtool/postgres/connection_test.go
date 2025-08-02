package postgres

import (
	"context"
	"os"
	"testing"
	"time"
)

// Requirements:
// - Set up Postgres database
// - Set environment variable
// "export TEST_DB_DSN=<dsn_of_postgres_db>"

// Reminder: There's no need to get env "TEST_DB_DSN" for every test function
// Just retrieve it one time with sync.Once
func TestGetPgxConn(t *testing.T) {
	testSingleConnCfg := SingleConnConfig{
		InitTimeout: 2 * time.Second,
	}

	dsn := os.Getenv("TEST_DB_DSN")
	conn, err := GetPgxConn(dsn, testSingleConnCfg)
	if err != nil {
		t.Errorf("failed to init connection to database: %s", err.Error())
		t.FailNow()
	}

	// Ping to database to test connection whether it's alive (rarely fail but just for test)
	pingCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := conn.Ping(pingCtx); err != nil {
		t.Errorf("failed to ping to database: %s", err.Error())
	}
}

func TestGetPgxPool(t *testing.T) {
	testPoolCfg := PoolConfig{
		InitTimeout: 2 * time.Second,
	}

	dsn := os.Getenv("TEST_DB_DSN")
	pool, err := GetPgxPool(dsn, testPoolCfg)
	if err != nil {
		t.Errorf("failed to init pool of connection to database: %s", err.Error())
		t.FailNow()
	}

	// Ping to database to test connection whether it's alive (rarely fail but just for test)
	pingCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		t.Errorf("failed to ping to database: %s", err.Error())
	}
}

func TestGetSqlxDB(t *testing.T) {
	testPoolCfg := PoolConfig{
		InitTimeout: 2 * time.Second,
	}

	dsn := os.Getenv("TEST_DB_DSN")
	pool, err := GetSqlxDB(dsn, testPoolCfg)
	if err != nil {
		t.Errorf("failed to init pool of connection to database: %s", err.Error())
		t.FailNow()
	}

	// Ping to database to test connection whether it's alive (rarely fail but just for test)
	pingCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := pool.PingContext(pingCtx); err != nil {
		t.Errorf("failed to ping to database: %s", err.Error())
	}
}
