package task

// Task service
type TaskService struct {
	repo TaskRepository
}

// Creates new TaskService
func NewTaskService(repo TaskRepository) TaskService {
	return TaskService{repo: repo}
}

// Get Task list
func (us *TaskService) GetTaskList() ([]Task, error) {
	users, err := us.repo.GetList()

	return users, err
}

// Save Task
func (us *TaskService) SaveTask(user *Task) error {
	return us.repo.Save(user)
}

// Delete Task by id
func (us *TaskService) DeleteTask(id int64) error {
	return us.repo.Delete(id)
}
