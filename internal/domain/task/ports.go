package task

import (
	"accounter/internal/domain/common"
	"accounter/pkg/utils"
	"bytes"
	"context"
)

// TaskRepository Task repository port
type TaskRepository interface {

	// Get list of Task
	GetList(ctx context.Context, params TaskParams) (Tasks, error)

	// Get one Task
	GetOne(ctx context.Context, id int64) (Task, error)

	// Create Task
	Create(ctx context.Context, task *Task) error

	// Update Task
	Update(ctx context.Context, task *Task) error

	// Delete one Task by id
	Delete(ctx context.Context, id int64) error

	// Execute operations in transaction
	WithTx(ctx context.Context, cb func(context.Context) error) error
}

// TaskRenderer Task renderer to varoius formats
type TaskRenderer interface {
	Render(format utils.FileFormat, tasks Tasks) (*bytes.Buffer, error)
}

// TaskEventsBus events bus for Task actions
type TaskEventsBus interface {
	// Trigger for create Task event
	OnCreate(ctx context.Context, userID int64, obj common.Entity) error

	// Trigger for update Task event
	OnUpdate(ctx context.Context, userID int64, old, new common.Entity) error

	// Trigger for delete Task event
	OnDelete(ctx context.Context, userID int64, obj common.Entity) error
}
