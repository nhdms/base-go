package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

func TestWithTransactionSqlx(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	testPoolCfg := PoolConfig{
		InitTimeout: 2 * time.Second,
	}
	db, err := GetSqlxDB(dsn, testPoolCfg)
	if err != nil {
		t.Errorf("failed to init pool of connection to database: %s", err.Error())
		t.FailNow()
	}

	// Create simple table just for test
	db.MustExec("create table if not exists example (field1 int, field2 text)")
	defer func() {
		db.MustExec("drop table example")
	}()

	testTxFunc := func(tx *sqlx.Tx) error {
		res, err := tx.Exec("insert into example values (123, 'hi')")
		if err != nil {
			return err
		}
		rowsAfffected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAfffected == 0 {
			return fmt.Errorf("failed to insert, there's no row affected")
		}

		res, err = tx.Exec("update example set field2='bye' where field1=123")
		if err != nil {
			return err
		}

		rowsAfffected, err = res.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAfffected == 0 {
			return fmt.Errorf("failed to update, there's no row affected")
		}

		return nil
	}

	if err := WithTransactionSqlx(context.Background(), db, testTxFunc); err != nil {
		t.Errorf("failed to process transaction: %s", err.Error())
	}
}

func TestWithTransactionPgx(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	testPoolCfg := PoolConfig{
		InitTimeout: 2 * time.Second,
	}
	pool, err := GetPgxPool(dsn, testPoolCfg)
	if err != nil {
		t.Errorf("failed to init pool of connection to database: %s", err.Error())
		t.FailNow()
	}

	// Create simple table just for test
	if _, err := pool.Exec(context.Background(), "create table if not exists example (field1 int, field2 text)"); err != nil {
		t.Errorf("failed to create table for test")
		t.FailNow()
	}
	defer func() {
		_, err := pool.Exec(context.Background(), "drop table example")
		if err != nil {
			slog.Error("failed to delete table after test", "error", err.Error())
		}
	}()

	testTxFunc := func(tx pgx.Tx) error {
		res, err := tx.Exec(context.Background(), "insert into example values (123, 'hi')")
		if err != nil {
			return err
		}

		if res.RowsAffected() == 0 {
			return fmt.Errorf("failed to insert, there's no row affected")
		}

		res, err = tx.Exec(context.Background(), "update example set field2='bye' where field1=123")
		if err != nil {
			return err
		}

		if res.RowsAffected() == 0 {
			return fmt.Errorf("failed to update, there's no row affected")
		}

		return nil
	}

	if err := WithTransactionPgx(context.Background(), pool, testTxFunc); err != nil {
		t.Errorf("failed to process transaction: %s", err.Error())
	}
}
