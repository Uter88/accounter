package postgres

import (
	"context"
)

// Base SQL repository
type baseRepository struct {
	client *SQLClient
}

// newBaseRepository creates new baseRepository
func newBaseRepository(client *SQLClient) *baseRepository {
	return &baseRepository{client: client}
}

// getContext creates timeout context with cancel func
func (r *baseRepository) getContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultOppTimeout)
}

// namedGet bind named query and exec GetContext
func (r *baseRepository) namedGet(ctx context.Context, query string, dest any, args any) error {
	db := r.client.GetExecutor(ctx)
	query, params, err := db.BindNamed(query, args)

	if err != nil {
		return err
	}

	return db.GetContext(ctx, dest, query, params...)
}

// namedSelect bind named query and exec SelectContext
func (r *baseRepository) namedSelect(ctx context.Context, query string, dest any, args any) error {
	db := r.client.GetExecutor(ctx)
	query, params, err := db.BindNamed(query, args)

	if err != nil {
		return err
	}

	return db.SelectContext(ctx, dest, query, params...)
}

// WithTx start transaction and execute callback
func (r *baseRepository) WithTx(ctx context.Context, cb func(context.Context) error) error {
	return r.client.BeginTx(ctx, cb)
}
