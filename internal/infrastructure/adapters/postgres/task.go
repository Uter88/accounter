package postgres

import (
	"accounter/internal/domain/task"
	"accounter/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Task repository
type taskRepository struct {
	*baseRepository
}

// NewTaskRepository creates new taskRepository
func NewTaskRepository(client *SQLClient) *taskRepository {
	return &taskRepository{
		baseRepository: newBaseRepository(client),
	}
}

// GetList get list of Task
func (r *taskRepository) GetList(ctx context.Context, params task.TaskParams) (task.Tasks, error) {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	result := make([]task.Task, 0)
	query := makeGetTaskQuery(params)
	err := r.namedSelect(ctx, query, &result, params)

	return result, err
}

// GetOne get one Task by id
func (r *taskRepository) GetOne(ctx context.Context, id int64) (t task.Task, err error) {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	db := r.client.GetExecutor(ctx)

	query := fmt.Sprintf("%s WHERE t.id = $1", getTaskQuery)
	err = db.GetContext(ctx, &t, query, id)

	if err == sql.ErrNoRows {
		return t, task.ErrTaskNotFound
	}

	return
}

// Create new Task
func (r *taskRepository) Create(ctx context.Context, t *task.Task) error {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	db := r.client.GetExecutor(ctx)

	if res, err := db.NamedExecContext(ctx, insertTaskQuery, t); err != nil {
		return err

	} else if id, _ := res.LastInsertId(); id != 0 {
		t.ID = id
	}

	return nil
}

// Update Task
func (r *taskRepository) Update(ctx context.Context, t *task.Task) error {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	db := r.client.GetExecutor(ctx)

	_, err := db.NamedExecContext(ctx, updateTaskQuery, t)

	return err
}

// Delete Task by id
func (r *taskRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	db := r.client.GetExecutor(ctx)
	_, err := db.ExecContext(ctx, deleteTaskQuery, id)

	return err
}

// makeGetTaskQuery build query from TaskParams
func makeGetTaskQuery(p task.TaskParams) string {
	query := getTaskQuery
	var conditions []string

	if !utils.IsEmpty(p.DateStart) && !utils.IsEmpty(p.DateEnd) {
		conditions = append(conditions, "t.date BETWEEN :date_start AND :date_end")
	}

	if !utils.IsEmpty(p.Status) {
		conditions = append(conditions, "t.status = :status")
	}

	if len(p.Users) > 0 {
		conditions = append(conditions, fmt.Sprintf("t.user_id IN (%s)", utils.Stringify(p.Users...)))
	}

	if len(conditions) > 0 {
		query += fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
	}

	if !utils.IsEmpty(p.OrderBy) {
		query += fmt.Sprintf(" ORDER BY t.%s", p.OrderBy)

		if p.OrderDesc {
			query += " DESC"
		}
	}

	if p.Limit > 0 {
		query += "LIMIT :skip, :limit"
	}

	return query
}

// Task queries
const (
	getTaskQuery = `
		SELECT
			t.id,
			t.user_id,
			CONCAT_WS(' ', u.surname, u.name) as user_label,
			u.price_per_hour,
			t.task_id,
			t.status,
			t.description,
			t.work_begin,
			t.work_end,
			t.date
		FROM tasks t
		JOIN users u ON u.id = t.user_id
	`
	deleteTaskQuery = `DELETE FROM tasks WHERE id = $1`
	updateTaskQuery = `
		UPDATE tasks SET
			user_id=:user_id,
			task_id=:task_id,
			status=:status,
			description=:description,
			work_begin=:work_begin,
			work_end=:work_end,
			date=:date
		WHERE id = :id
	`
	insertTaskQuery = `
		INSERT INTO tasks (user_id, task_id, status, description, work_begin, work_end, date)
		VALUES (:user_id, :task_id, :status, :description, :work_begin, :work_end, :date)
		RETURNING id;
	`
)
