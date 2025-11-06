package task

// Task model
type Task struct {
	ID int64 `db:"id,omitempty" json:"id"`
}
