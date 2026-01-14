package store

import (
	"accounter/internal/domain/common"
	"accounter/internal/domain/task"
	"accounter/pkg/utils"
	"fmt"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type tasksStore struct {
	*baseStore

	tasks   task.Tasks
	params  *task.TaskParams
	loading bool
}

func (s *tasksStore) RequestTasks(ctx app.Context) error {
	s.setLoading(true)
	defer s.setLoading(false)

	resp, errResp, err := newRequest[task.Tasks](*s.baseStore).
		Path("tasks/list").
		Params(s.params.Encode()).
		Do()

	if err != nil {
		return fmt.Errorf("error request tasks, status: %s, error: %s", err.Error(), errResp.Error)
	}

	s.tasks = resp.Data

	ctx.NewActionWithValue("setTasks", s.tasks)

	return nil
}

func (s *tasksStore) SaveTask(ctx app.Context, t task.Task) (task.Task, error) {
	resp, errResp, err := newRequest[task.Task](*s.baseStore).
		Path("tasks/save").
		Method(http.MethodPost).
		Data(utils.ToJSON(t)).
		Do()

	if err != nil {
		return t, fmt.Errorf("error save task, status: %s, error: %s", err.Error(), errResp.Error)
	}

	result := resp.Data

	return result, s.RequestTasks(ctx)
}

func (s *tasksStore) RemoveTask(ctx app.Context, t task.Task) error {
	_, errResp, err := newRequest[string](*s.baseStore).
		Path(fmt.Sprintf("tasks/delete/%d", t.ID)).
		Method(http.MethodDelete).
		Do()

	if err != nil {
		return fmt.Errorf("error save task, status: %s, error: %s", err.Error(), errResp.Error)
	}

	return s.RequestTasks(ctx)
}

func (s *tasksStore) ExportTasks(format utils.FileFormat) string {
	request := newRequest[string](*s.baseStore).
		Path(fmt.Sprintf("tasks/export/%s", format)).
		Params(s.params.Encode()).
		Param("token", s.GetUser().Tokens.AccessToken).
		Method(http.MethodGet)

	return request.GetURL()
}

func (s *tasksStore) GetTasks() task.Tasks {
	return s.tasks
}

func (s *tasksStore) SetTaskParams(p *task.TaskParams) {
	s.params = p

	s.ws.SendMessage(common.WsMessageParams, p)
}

func (s *tasksStore) GetTaskParams() *task.TaskParams {
	return s.params
}

func (s *tasksStore) GetTasksLoading() bool {
	return s.loading
}

func (s *tasksStore) setLoading(v bool) {
	s.loading = v
}
