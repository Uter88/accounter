package adapter_sql

import (
	"accounter/domain/task"
	"context"
)

// Task repository
type taskRepository struct {
	baseRepository
}

// Creates new taskRepository
func NewTaskRepository(ctx context.Context, client SQLClient) *taskRepository {
	return &taskRepository{
		baseRepository: newBaseRepository(ctx, client),
	}
}

// Get list of Task
func (r *taskRepository) GetList() ([]task.Task, error) {
	return nil, nil
}

// Get one Task by id
func (r *taskRepository) GetOne(id int64) (task.Task, error) {
	return task.Task{}, nil
}

// Save Task
func (r *taskRepository) Save(t *task.Task) error {
	return nil
}

// Delete Task by id
func (r *taskRepository) Delete(id int64) error {
	return nil
}
