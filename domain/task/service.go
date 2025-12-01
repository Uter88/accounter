package task

import (
	"accounter/pkg/tools"
	"bytes"
	"errors"
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
	result, err := ts.repo.GetList(params)

	return result, err
}

// Save Task
func (ts *TaskService) SaveTask(task *Task) error {
	if tools.IsEmpty(task.ID) {
		return ts.repo.Insert(task)
	}

	return ts.repo.Update(task)
}

// Delete Task by id
func (ts *TaskService) DeleteTask(id int64) error {
	return ts.repo.Delete(id)
}

// Export Tasks by specified format
func (ts *TaskService) ExportTasks(tasks Tasks, format tools.FileFormat) (*bytes.Buffer, error) {
	data := tasksToExportData(tasks)

	switch format {
	case tools.FileFormatCSV:
		return data.convertToCSV()

	case tools.FileFormatDocX:
		return data.convertToDocX()

	case tools.FileFormatXLSX:
		return data.convertToXLSX()

	case tools.FileFormatPDF:
		return data.convertToPDF()

	case tools.FileFormatJSON:
		return data.convertToJSON()

	case tools.FileFormatHTML:
		return data.convertToHTML()

	default:
		return nil, errors.New("unexpected file format")
	}
}
