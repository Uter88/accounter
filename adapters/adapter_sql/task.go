package adapter_sql

import (
	"accounter/domain/task"
	"accounter/pkg/tools"
	"context"
	"fmt"
	"strings"
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
func (r *taskRepository) GetList(params *task.TaskParams) ([]task.Task, error) {
	ctx, cancel := r.getContext()
	defer cancel()

	result := make([]task.Task, 0)
	query := makeGetTaskQuery(params)
	err := r.namedSelect(ctx, query, &result, params)

	return result, err
}

// Get one Task by id
func (r *taskRepository) GetOne(id int64) (t task.Task, err error) {
	ctx, cancel := r.getContext()
	defer cancel()

	query := fmt.Sprintf("%s WHERE t.id = $1", getTaskQuery)
	err = r.db().GetContext(ctx, &t, query)

	return
}

// Save Task
func (r *taskRepository) Insert(t *task.Task) error {
	ctx, cancel := r.getContext()
	defer cancel()

	if res, err := r.db().NamedExecContext(ctx, insertTaskQuery, t); err != nil {
		return err

	} else if id, _ := res.LastInsertId(); id != 0 {
		t.ID = id
	}

	return nil
}

// Update Task
func (r *taskRepository) Update(t *task.Task) error {
	ctx, cancel := r.getContext()
	defer cancel()

	_, err := r.db().NamedExecContext(ctx, updateTaskQuery, t)

	return err
}

// Delete Task by id
func (r *taskRepository) Delete(id int64) error {
	ctx, cancel := r.getContext()
	defer cancel()

	_, err := r.db().ExecContext(ctx, deleteTaskQuery, id)

	return err
}

func makeGetTaskQuery(p *task.TaskParams) string {
	query := getTaskQuery
	var conditions []string

	if !tools.IsEmpty(p.DateStart) && !tools.IsEmpty(p.DateEnd) {
		conditions = append(conditions, "t.date BETWEEN :date_start AND :date_end")
	}

	if !tools.IsEmpty(p.Status) {
		conditions = append(conditions, "t.status = :status")
	}

	if len(p.Users) > 0 {
		conditions = append(conditions, fmt.Sprintf("t.user_id IN (%s)", tools.Stringify(p.Users...)))
	}

	if len(conditions) > 0 {
		query += fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
	}

	if !tools.IsEmpty(p.OrderBy) {
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
	`
)
