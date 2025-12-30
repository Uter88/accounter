package adapter_sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	// Key of context value for transaction
	txKey = "tx"

	// Default operation and connection timeouts
	defaultOppTimeout        = time.Minute * 5
	defaultConnectionTimeout = time.Second * 20
)

// SQL client
type SQLClient struct {
	driver string
	dsn    string
	db     *sqlx.DB
}

// Common sqlx executor
type executor interface {
	sqlx.Ext
	sqlx.ExtContext
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}

// Creates new SQLClient
func NewSQLClient(driver, dsn string) SQLClient {
	return SQLClient{driver: driver, dsn: dsn}
}

// Connect to database
func (c *SQLClient) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultConnectionTimeout)
	defer cancel()

	if conn, err := sqlx.Open(c.driver, c.dsn); err != nil {
		return fmt.Errorf("could not create %s DB connection pool: %s", c.driver, err.Error())

	} else if err = conn.PingContext(ctx); err != nil {
		return fmt.Errorf("error ping %s DB connection: %s", c.driver, err.Error())

	} else {
		c.db = conn
	}

	return nil
}

// Disconnect from database
func (c SQLClient) Disconnect() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

// GetExecutor returns transaction from context or common connection
func (c SQLClient) GetExecutor(ctx context.Context) executor {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx
	}

	return c.db
}

// BeginTx start new transaction and put it to context
// in case of error - will auto rollback transactions
func (c SQLClient) BeginTx(ctx context.Context, cb func(context.Context) error) error {
	tx, err := c.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	ctx = context.WithValue(ctx, txKey, tx)

	if err = cb(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
