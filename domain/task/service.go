package task

import (
	"accounter/pkg/tools"
	"bytes"
)

// Task service
type TaskService struct {
	repo TaskRepository
}

// Creates new TaskService
func NewTaskService(repo TaskRepository) TaskService {
	return TaskService{repo: repo}
}

// Get Task list
func (ts *TaskService) GetTaskList(params *TaskParams) ([]Task, error) {
	users, err := ts.repo.GetList(params)

	return users, err
}

// Save Task
func (ts *TaskService) SaveTask(user *Task) error {
	return ts.repo.Save(user)
}

// Delete Task by id
func (ts *TaskService) DeleteTask(id int64) error {
	return ts.repo.Delete(id)
}

func (ts *TaskService) ExportTasks(tasks Tasks, format tools.FileFormat) (*bytes.Buffer, error) {
	return exportTasks(tasks, format)
}
