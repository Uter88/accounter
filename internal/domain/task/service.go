package task

import (
	"accounter/pkg/tools"
	"bytes"
	"context"
)

// Task service
type TaskService struct {
	repo     TaskRepository
	renderer TaskRenderer
}

// Creates new TaskService
func NewTaskService(repo TaskRepository, renderer TaskRenderer) TaskService {
	return TaskService{repo: repo, renderer: renderer}
}

// Get Task list
func (ts *TaskService) GetTaskList(ctx context.Context, params *TaskParams) ([]Task, error) {
	result, err := ts.repo.GetList(ctx, params)

	return result, err
}

// Save Task
func (ts *TaskService) SaveTask(ctx context.Context, task *Task) error {
	if tools.IsEmpty(task.ID) {
		return ts.repo.Insert(ctx, task)
	}

	return ts.repo.Update(ctx, task)
}

// Delete Task by id
func (ts *TaskService) DeleteTask(ctx context.Context, id int64) error {
	return ts.repo.Delete(ctx, id)
}

// Export Tasks by specified format
func (ts *TaskService) ExportTasks(tasks Tasks, format tools.FileFormat) (*bytes.Buffer, error) {
	return ts.renderer.Render(format, tasks)
}
