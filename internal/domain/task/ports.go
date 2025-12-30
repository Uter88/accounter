package task

import (
	"accounter/pkg/tools"
	"bytes"
	"context"
)

// TaskRepository Task repository port
type TaskRepository interface {

	// Get list of Task
	GetList(ctx context.Context, params *TaskParams) ([]Task, error)

	// Get one Task
	GetOne(ctx context.Context, id int64) (Task, error)

	// Create Task
	Insert(ctx context.Context, task *Task) error

	// Update Task
	Update(ctx context.Context, task *Task) error

	// Delete one Task by id
	Delete(ctx context.Context, id int64) error

	// Execute operations in transaction
	WithTx(ctx context.Context, cb func(context.Context) error) error
}

// TaskRenderer Task renderer to varoius formats
type TaskRenderer interface {
	Render(format tools.FileFormat, tasks Tasks) (*bytes.Buffer, error)
}
