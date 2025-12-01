package task

// Task repository port
type TaskRepository interface {

	// Get list of Task
	GetList(params *TaskParams) ([]Task, error)

	// Get one Task
	GetOne(id int64) (Task, error)

	// Create Task
	Insert(*Task) error

	// Update Task
	Update(*Task) error

	// Delete one Task by id
	Delete(id int64) error
}
