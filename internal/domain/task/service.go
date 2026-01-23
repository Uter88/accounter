package task

import (
	"accounter/internal/domain/common"
	"accounter/pkg/utils"
	"bytes"
	"context"
)

// Task service
type TaskService struct {
	repo     TaskRepository
	renderer TaskRenderer
	eventBus TaskEventsBus
}

// Creates new TaskService
func NewTaskService(repo TaskRepository, renderer TaskRenderer, eventBus TaskEventsBus) *TaskService {
	return &TaskService{
		repo:     repo,
		renderer: renderer,
		eventBus: eventBus,
	}
}

// Get Task list
func (ts *TaskService) GetTaskList(ctx context.Context, params common.RequestParams) (Tasks, error) {
	return ts.repo.GetList(ctx, params)
}

// GetTask get one Task by id
func (ts *TaskService) GetTask(ctx context.Context, id int64) (Task, error) {
	return ts.repo.GetOne(ctx, id)
}

// Save create or update Task
func (ts *TaskService) SaveTask(ctx common.Context, task *Task) error {
	if utils.IsEmpty(task.ID) {
		return ts.createTask(ctx, task)
	}

	return ts.updateTask(ctx, task)
}

// updateTask updates exitance Task
func (ts *TaskService) updateTask(ctx common.Context, task *Task) error {
	userID := ctx.GetID()
	oldTask, err := ts.GetTask(ctx, task.ID)

	if err != nil {
		return err
	}

	return ts.repo.WithTx(ctx, func(ctx context.Context) error {
		if err = ts.repo.Update(ctx, task); err != nil {
			return ErrTaskUpdate.WithErr(err)
		}

		if err := ts.eventBus.OnUpdate(ctx, userID, oldTask, *task); err != nil {
			return ErrLogTask.WithErr(err)
		}

		return nil
	})
}

// createTask creates new Task
func (ts *TaskService) createTask(ctx common.Context, task *Task) error {
	userID := ctx.GetID()

	return ts.repo.WithTx(ctx, func(ctx context.Context) error {
		if err := ts.repo.Create(ctx, task); err != nil {
			return ErrTaskCreate.WithErr(err)
		}

		if err := ts.eventBus.OnCreate(ctx, userID, task); err != nil {
			return ErrLogTask.WithErr(err)
		}

		return nil
	})
}

// Delete Task
func (ts *TaskService) DeleteTask(ctx common.Context, id int64) error {
	userID := ctx.GetID()

	task, err := ts.GetTask(ctx, id)

	if err != nil {
		return err
	}

	return ts.repo.WithTx(ctx, func(ctx context.Context) error {
		if err := ts.repo.Delete(ctx, task.ID); err != nil {
			return ErrTaskDelete.WithErr(err)
		}

		if err := ts.eventBus.OnDelete(ctx, userID, task); err != nil {
			return ErrLogTask.WithErr(err)
		}

		return nil
	})
}

// Export Tasks by specified format
func (ts *TaskService) ExportTasks(tasks Tasks, format utils.FileFormat) (*bytes.Buffer, error) {
	return ts.renderer.Render(format, tasks)
}
