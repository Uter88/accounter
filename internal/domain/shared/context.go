package shared

import (
	"context"
)

// Context with ID of current User
type Context interface {
	context.Context

	// ID of current User
	GetID() int64
}
