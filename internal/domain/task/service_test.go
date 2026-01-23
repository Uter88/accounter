package task

import (
	"accounter/internal/domain/common"
	"accounter/pkg/utils"
	"bytes"
	"context"
	"errors"
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// Service testing suite
type testServiceSuite struct {
	ctx common.Context
	suite.Suite
	sampleService *TaskService
	fixedTime     time.Time
	eventsChan    chan common.Entity
}

// Mock for TaskRepository
type mockTaskRepostory struct {
	store Tasks
}

// GetList mock of returning Tasks
func (m *mockTaskRepostory) GetList(ctx context.Context, p common.RequestParams) (Tasks, error) {
	return m.store, nil
}

// Create mock of Task creation
func (m *mockTaskRepostory) Create(ctx context.Context, task *Task) error {
	if err, ok := ctx.Value("oppError").(error); ok {
		return err
	}

	task.ID = rand.Int63()
	m.store = append(m.store, *task)
	return nil
}

// Update mock of Task updation
func (m *mockTaskRepostory) Update(ctx context.Context, task *Task) error {
	if err, ok := ctx.Value("oppError").(error); ok {
		return err
	}

	for i, t := range m.store {
		if t.ID == task.ID {
			m.store[i] = *task
			break
		}
	}

	return nil
}

// Delete mock of Task deletion
func (m *mockTaskRepostory) Delete(ctx context.Context, id int64) error {
	if err, ok := ctx.Value("oppError").(error); ok {
		return err
	}

	m.store = slices.DeleteFunc(m.store, func(t Task) bool { return t.ID == id })
	return nil
}

// GetOne mock of getting single Task
func (m *mockTaskRepostory) GetOne(ctx context.Context, id int64) (Task, error) {
	for _, t := range m.store {
		if t.ID == id {
			return t, nil
		}
	}

	return Task{}, ErrTaskNotFound
}

// WithTx mock of begin Task transaction
func (m *mockTaskRepostory) WithTx(ctx context.Context, cb func(context.Context) error) error {
	return cb(ctx)
}

// mockTaskRenderer mock for Task renderer
type mockTaskRenderer struct{}

// Render mock of Tasks rendering
func (m *mockTaskRenderer) Render(format utils.FileFormat, tasks Tasks) (*bytes.Buffer, error) {
	return bytes.NewBuffer(nil), nil
}

// mockEventBus mock of Task event bus
type mockEventBus struct {
	queue chan common.Entity
}

// OnCreate mock of create event
func (m *mockEventBus) OnCreate(ctx context.Context, userID int64, obj common.Entity) error {
	if err, ok := ctx.Value("evError").(error); ok {
		return err
	}

	m.queue <- obj
	return nil
}

// OnDelete mock of delete event
func (m *mockEventBus) OnDelete(ctx context.Context, userID int64, obj common.Entity) error {
	if err, ok := ctx.Value("evError").(error); ok {
		return err
	}

	m.queue <- obj
	return nil
}

// OnUpdate mock of update event
func (m *mockEventBus) OnUpdate(ctx context.Context, userID int64, old, new common.Entity) error {
	if err, ok := ctx.Value("evError").(error); ok {
		return err
	}

	m.queue <- new
	return nil
}

// TestTaskSuite run Task testing
func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(testServiceSuite))
}

// SetupTest testing preparations
func (suite *testServiceSuite) SetupTest() {
	suite.fixedTime = time.Now()
	suite.ctx = common.NewTestContext(context.Background(), 1)
	suite.eventsChan = make(chan common.Entity, 1)

	repo := &mockTaskRepostory{
		store: Tasks{
			{ID: 1, Description: "Test task 1", Date: suite.fixedTime.Unix()},
			{ID: 2, Description: "Test task 2", Date: suite.fixedTime.Unix()},
		},
	}
	renderer := new(mockTaskRenderer)
	eventBus := &mockEventBus{
		queue: suite.eventsChan,
	}

	suite.sampleService = NewTaskService(repo, renderer, eventBus)
}

// TestGetTaskList testing of Tasks getting
func (suite *testServiceSuite) TestGetTaskList() {
	params := common.NewRequestParams(suite.fixedTime)

	result, err := suite.sampleService.GetTaskList(suite.ctx, params)

	suite.Nil(err, "Unexpected error")
	suite.Equal(len(result), 2, "Unexpected tasks length")
}

// TestGetTask testing of getting single Task
func (suite *testServiceSuite) TestGetTask() {
	id := int64(1)
	task, err := suite.sampleService.GetTask(suite.ctx, id)

	suite.Nil(err, "Unexpected error")
	suite.Equal(task.ID, id, "Invalid id")
}

// TestSaveTask test of saving Task
func (suite *testServiceSuite) TestSaveTask() {
	existTask := NewTask(suite.fixedTime)
	existTask.ID = 1
	newTask := NewTask(suite.fixedTime)

	err := suite.sampleService.SaveTask(suite.ctx, &existTask)
	suite.Nil(err, "Unexpected error")
	suite.Equal(len(suite.eventsChan), 1, "Empty events queue")
	ev := <-suite.eventsChan
	suite.Equal(existTask.ID, ev.GetID(), "Incorrect task id")

	err = suite.sampleService.SaveTask(suite.ctx, &newTask)
	suite.Nil(err, "Unexpected error")
	suite.Equal(len(suite.eventsChan), 1, "Empty events queue")
	ev = <-suite.eventsChan
	suite.NotEqual(existTask.ID, ev.GetID(), "Incorrect task id")
}

// TestCreateTask testing of creating Task
func (suite *testServiceSuite) TestCreateTask() {
	task := NewTask(suite.fixedTime)

	err := suite.sampleService.createTask(suite.ctx, &task)
	suite.Nil(err, "Unexpected error")

	suite.Equal(len(suite.eventsChan), 1, "Empty events queue")
	ev := <-suite.eventsChan
	suite.Equal(task.ID, ev.GetID(), "Incorrect task id")

	task.ID = 0
	errCtx := context.WithValue(suite.ctx, "oppError", errors.New("some error"))
	ctx := common.NewTestContext(errCtx, 1)
	err = suite.sampleService.createTask(ctx, &task)
	suite.True(common.IsDomainError(err, ErrTaskCreate.Code))

	errCtx = context.WithValue(suite.ctx, "evError", errors.New("some error"))
	ctx = common.NewTestContext(errCtx, 1)
	err = suite.sampleService.createTask(ctx, &task)
	suite.True(common.IsDomainError(err, ErrLogTask.Code))
}

// TestUpdateTask testing of updating Task
func (suite *testServiceSuite) TestUpdateTask() {
	task := NewTask(suite.fixedTime)
	task.ID = 1

	err := suite.sampleService.updateTask(suite.ctx, &task)
	suite.Nil(err, "Unexpected error")

	suite.Equal(len(suite.eventsChan), 1, "Empty events queue")
	ev := <-suite.eventsChan
	suite.Equal(task.ID, ev.GetID(), "Incorrect task id")

	task.ID = 0
	err = suite.sampleService.updateTask(suite.ctx, &task)
	suite.ErrorIs(err, ErrTaskNotFound)

	task.ID = 1

	errCtx := context.WithValue(suite.ctx, "oppError", errors.New("some error"))
	ctx := common.NewTestContext(errCtx, 1)
	err = suite.sampleService.updateTask(ctx, &task)
	suite.True(common.IsDomainError(err, ErrTaskUpdate.Code))

	errCtx = context.WithValue(suite.ctx, "evError", errors.New("some error"))
	ctx = common.NewTestContext(errCtx, 1)
	err = suite.sampleService.updateTask(ctx, &task)
	suite.True(common.IsDomainError(err, ErrLogTask.Code))
}

// TestDeleteTask testing of Task deleting
func (suite *testServiceSuite) TestDeleteTask() {
	id := int64(1)

	err := suite.sampleService.DeleteTask(suite.ctx, id)
	suite.Nil(err, "Unexpected error")

	suite.Equal(len(suite.eventsChan), 1, "Empty events queue")
	ev := <-suite.eventsChan
	suite.Equal(id, ev.GetID(), "Incorrect task id")

	err = suite.sampleService.DeleteTask(suite.ctx, id)
	suite.ErrorIs(err, ErrTaskNotFound)

	errCtx := context.WithValue(suite.ctx, "oppError", errors.New("some error"))
	ctx := common.NewTestContext(errCtx, 1)
	err = suite.sampleService.DeleteTask(ctx, 2)
	suite.True(common.IsDomainError(err, ErrTaskDelete.Code))

	errCtx = context.WithValue(suite.ctx, "evError", errors.New("some error"))
	ctx = common.NewTestContext(errCtx, 1)
	err = suite.sampleService.DeleteTask(ctx, 2)
	suite.True(common.IsDomainError(err, ErrLogTask.Code))
}

// TestExportTasks testing of Tasks export
func (suite *testServiceSuite) TestExportTasks() {
	var tasks Tasks

	result, err := suite.sampleService.ExportTasks(tasks, utils.FileFormatJSON)
	suite.Nil(err, "Unexpected error")
	suite.IsType(result, &bytes.Buffer{})
}
