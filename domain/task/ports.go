package task

// Task repository port
type TaskRepository interface {

	// Get list of Task
	GetList() ([]Task, error)

	// Get one Task
	GetOne(id int64) (Task, error)

	// Create/update Task
	Save(*Task) error

	// Delete one Task by id
	Delete(id int64) error
}
