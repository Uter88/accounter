package adapter_sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Base SQL repository
type baseRepository struct {
	client SQLClient
}

// Creates new baseRepository
func newBaseRepository(client SQLClient) baseRepository {
	return baseRepository{client: client}
}

// Creates timeout context with cancel func
func (r *baseRepository) getContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, time.Minute*5)
}

func (r *baseRepository) namedGet(ctx context.Context, query string, dest any, args any) error {
	query, params, err := r.client.db.BindNamed(query, args)

	if err != nil {
		return err
	}

	return r.db().GetContext(ctx, dest, query, params...)
}

func (r *baseRepository) namedSelect(ctx context.Context, query string, dest any, args any) error {
	query, params, err := r.client.db.BindNamed(query, args)

	if err != nil {
		return err
	}

	return r.db().SelectContext(ctx, dest, query, params...)
}

// Returns SQLX database instance
func (r *baseRepository) db() *sqlx.DB {
	return r.client.db
}

// SQL client
type SQLClient struct {
	driver string
	dsn    string
	db     *sqlx.DB
}

// Creates new SQLClient
func NewSQLClient(driver, dsn string) SQLClient {
	return SQLClient{driver: driver, dsn: dsn}
}

// Connect to database
func (c *SQLClient) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
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
func (c *SQLClient) Disconnect() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

func (c *SQLClient) DB() *sqlx.DB {
	return c.db
}
