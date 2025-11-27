package task

import (
	"accounter/pkg/tools"
	"fmt"
	"time"
)

// Task slice model
type Tasks []Task

func (tasks Tasks) Empty() bool { return tasks.Len() == 0 }
func (tasks Tasks) Len() int    { return len(tasks) }

func (tasks Tasks) GetPrice() (price float32) {
	for i := range tasks {
		price += tasks[i].GetPrice()
	}

	return tools.ToFixed(price, 2)
}

func (tasks Tasks) GroupByUsers() map[string]Tasks {
	result := make(map[string]Tasks)

	for _, t := range tasks {
		result[t.UserLabel] = append(result[t.UserLabel], t)
	}

	return result
}

func (tasks Tasks) GroupByDates() map[int64]Tasks {
	result := make(map[int64]Tasks)

	for _, t := range tasks {
		result[t.Date] = append(result[t.Date], t)
	}

	return result
}

// Task model
type Task struct {
	ID           int64   `db:"id,omitempty" json:"id,omitempty"`
	UserID       int64   `db:"user_id" json:"user_id"`
	UserLabel    string  `db:"user_label" json:"user_label"`
	PricePerHour float32 `db:"price_per_hour" json:"price_per_hour"`
	TaskID       string  `db:"task_id" json:"task_id"`
	Status       string  `db:"status" json:"status"`
	Description  string  `db:"description" json:"description"`
	WorkBegin    int64   `db:"work_begin" json:"work_begin"`
	WorkEnd      int64   `db:"work_end" json:"work_end"`
	Date         int64   `db:"date" json:"date"`
}

func (t Task) ToMap() tools.Data {
	result := make(tools.Data)

	result["id"] = tools.ValOrNil(t.ID)
	result["user_id"] = t.UserID
	result["description"] = t.Description
	result["task_id"] = t.TaskID
	result["status"] = t.Status
	result["date"] = t.Date
	result["work_begin"] = t.WorkBegin
	result["work_end"] = t.WorkEnd

	return result
}

func NewTask() Task {
	now := time.Now()

	return Task{
		Status:    "completed",
		WorkBegin: time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location()).Unix(),
		WorkEnd:   time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location()).Unix(),
		Date:      now.Unix(),
	}
}

func (t *Task) IsValid() bool {
	if tools.IsSomeEmpty(t.UserID, t.Date, t.WorkBegin, t.WorkEnd) {
		return false
	}

	if tools.IsSomeEmpty(t.Description, t.Status) {
		return false
	}

	return true
}

func (t *Task) SetDate(dt int64) *Task {
	t.Date = dt
	d := time.Unix(dt, 0)
	b := time.Unix(t.WorkBegin, 0)
	e := time.Unix(t.WorkEnd, 0)

	t.WorkBegin = time.Date(d.Year(), d.Month(), d.Day(), b.Hour(), b.Minute(), 0, 0, d.Location()).Unix()
	t.WorkEnd = time.Date(d.Year(), d.Month(), d.Day(), e.Hour(), e.Minute(), 0, 0, d.Location()).Unix()

	return t
}

func (t *Task) GetPrice() float32 {
	hours := float32(t.GetDuration().Hours())

	return tools.ToFixed(hours*t.PricePerHour, 2)
}

func (t *Task) FormatPrice() string {
	price := t.GetPrice()

	return fmt.Sprintf("%.2f", price)
}

func (t *Task) FormatDate() string {
	return time.Unix(t.Date, 0).Format("02.01.2006")
}

func (t *Task) FormatWorkBegin() string {
	return time.Unix(t.WorkBegin, 0).Format("15:04")
}

func (t *Task) FormatWorkEnd() string {
	return time.Unix(t.WorkEnd, 0).Format("15:04")
}

func (t *Task) FormatDuration(withSeconds bool) string {
	return tools.FormatDuration(t.GetDuration(), withSeconds)
}

func (t *Task) GetDuration() time.Duration {
	return time.Unix(t.WorkEnd, 0).Sub(time.Unix(t.WorkBegin, 0))
}
