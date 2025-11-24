package store

import (
	"accounter/domain/task"
	"accounter/pkg/tools"
	"fmt"
	"net/http"
)

type tasksStore struct {
	*baseStore

	tasks  task.Tasks
	params *task.TaskParams
}

func (s *tasksStore) RequestTasks() error {
	resp, errResp, err := newRequest[task.Tasks](*s.baseStore).
		Path("tasks/list").
		Params(s.params.Encode()).
		Do()

	if err != nil {
		return fmt.Errorf("error request tasks, status: %s, error: %s", err.Error(), errResp.Error)
	}

	s.tasks = resp.Data

	return nil
}

func (s *tasksStore) SaveTask(t task.Task) (task.Task, error) {
	resp, errResp, err := newRequest[task.Task](*s.baseStore).
		Path("tasks/save").
		Method(http.MethodPost).
		Data(tools.ToJSON(t)).
		Do()

	if err != nil {
		return t, fmt.Errorf("error save task, status: %s, error: %s", err.Error(), errResp.Error)
	}

	result := resp.Data

	return result, s.RequestTasks()
}

func (s *tasksStore) RemoveTask(t task.Task) error {
	_, errResp, err := newRequest[string](*s.baseStore).
		Path(fmt.Sprintf("tasks/delete/%d", t.ID)).
		Method(http.MethodDelete).
		Do()

	if err != nil {
		return fmt.Errorf("error save task, status: %s, error: %s", err.Error(), errResp.Error)
	}

	return s.RequestTasks()
}

func (s *tasksStore) GetTasks() task.Tasks {
	return s.tasks
}

func (s *tasksStore) SetTaskParams(p *task.TaskParams) {
	s.params = p
}

func (s *tasksStore) GetTaskParams() *task.TaskParams {
	return s.params
}
