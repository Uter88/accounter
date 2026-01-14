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

// GetPrice get total price of Tasks
func (tasks Tasks) GetPrice() (price float32) {
	for i := range tasks {
		price += tasks[i].GetPrice()
	}

	return tools.ToFixed(price, 2)
}

// GetDuration get total duration of Tasks
func (tasks Tasks) GetDuration() (duration time.Duration) {
	for i := range tasks {
		duration += tasks[i].GetDuration()
	}

	return
}

// GroupByUsers group Tasks by User id
func (tasks Tasks) GroupByUsers() *tools.OrderedMap[string, Tasks] {
	result := tools.NewOrderedMap[string, Tasks]()

	for _, t := range tasks {
		items, _ := result.Get(t.UserLabel)
		items = append(items, t)
		result.Set(t.UserLabel, items)
	}

	result.Sort()

	return result
}

// GroupByDates group Tasks by date
func (tasks Tasks) GroupByDates() *tools.OrderedMap[int64, Tasks] {
	result := tools.NewOrderedMap[int64, Tasks]()

	for _, t := range tasks {
		items, _ := result.Get(t.Date)
		items = append(items, t)
		result.Set(t.Date, items)
	}

	result.Sort()

	return result
}

// Fields for detect differences on update Event
var ComparationFields = [...]string{
	"UserID", "TaskID", "Status", "Description", "WorkBegin", "WorkEnd", "Date",
}

// Task entity
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

// NewTask creates new Task
func NewTask() Task {
	now := time.Now()

	return Task{
		Status:    "completed",
		WorkBegin: time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location()).Unix(),
		WorkEnd:   time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location()).Unix(),
		Date:      now.Unix(),
	}
}

// SetDate set Task date, begin and end of work
func (t *Task) SetDate(dt int64) *Task {
	t.Date = dt
	t.WorkBegin = tools.SetDateTs(dt, t.WorkBegin)
	t.WorkEnd = tools.SetDateTs(dt, t.WorkEnd)
	return t
}

// IsValid Task complete validation
func (t Task) IsValid() bool {
	if tools.IsSomeEmpty(t.UserID, t.Date, t.WorkBegin, t.WorkEnd) {
		return false
	}

	if tools.IsSomeEmpty(t.Description, t.Status) {
		return false
	}

	return true
}

// GetPrice get Task price by hours*price per hour
func (t Task) GetPrice() float32 {
	hours := float32(t.GetDuration().Hours())
	return tools.ToFixed(hours*t.PricePerHour, 2)
}

// FormatPrice string representation of price
func (t Task) FormatPrice() string {
	price := t.GetPrice()

	return fmt.Sprintf("%.2f", price)
}

// FormatDate string representation of date
func (t Task) FormatDate() string {
	return time.Unix(t.Date, 0).Format("02.01.2006")
}

// FormatWorkBegin string representation of work begin
func (t Task) FormatWorkBegin() string {
	return time.Unix(t.WorkBegin, 0).Format("15:04")
}

// FormatWorkEnd string representation of work end
func (t Task) FormatWorkEnd() string {
	return time.Unix(t.WorkEnd, 0).Format("15:04")
}

// FormatDuration string representation of time duration
func (t Task) FormatDuration(withSeconds bool) string {
	return tools.FormatDuration(t.GetDuration(), withSeconds)
}

// GetDuration get duration between work begin and work end dates
func (t Task) GetDuration() time.Duration {
	return time.Unix(t.WorkEnd, 0).Sub(time.Unix(t.WorkBegin, 0))
}

// GetID return Task id
func (t Task) GetID() int64 {
	return t.ID
}

// GetType return type of Task object
func (t Task) GetType() string {
	return "task"
}
