package common

import (
	"context"
)

// Context with ID of current User
type Context interface {
	context.Context

	// ID of current User
	GetID() int64
}

type testContext struct {
	context.Context
	id int64
}

func (t *testContext) GetID() int64 {
	return t.id
}

func NewTestContext(ctx context.Context, id int64) Context {
	return &testContext{Context: ctx, id: id}
}
