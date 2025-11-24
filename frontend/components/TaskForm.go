package components

import (
	"accounter/domain/task"
	"accounter/frontend/common"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type TaskForm struct {
	app.Compo
	common.BaseComponent

	visible bool
	data    task.Task
}

func NewTaskForm(ctx common.AppContext) *TaskForm {
	return &TaskForm{
		BaseComponent: common.NewBaseComponent(ctx),
		data:          task.NewTask(),
	}
}

func (uf *TaskForm) OnMount(ctx app.Context) {
	ctx.Handle("setTask", func(ctx app.Context, a app.Action) {
		task := a.Value.(task.Task)
		uf.data = task
		uf.Open()
	})

	ctx.Handle("resetTask", func(ctx app.Context, a app.Action) {
		uf.data = task.NewTask()
		uf.Hide()
	})
}

func (tf *TaskForm) Open() *TaskForm {
	tf.visible = true
	return tf
}

func (tf *TaskForm) Hide() *TaskForm {
	tf.visible = false
	return tf
}

func (tf *TaskForm) El() app.UI {
	return app.Form().
		Class("d-flex flex-column").
		Body(
			NewInputField[string]().
				Label("Description").
				Value(&tf.data.Description).
				WrapClass("mt-4").
				Clearable(true).
				Required(true).
				Autofocus(true).
				ID("description-field"),

			NewInputField[string]().
				Label("ID").
				Value(&tf.data.TaskID).
				WrapClass("mt-4").
				Clearable(true).
				ID("task-id-field"),

			NewSelectField[string]().
				Label("Status").
				Value(&tf.data.Status).
				Options([]SelectOption[string]{
					{
						Label: "In progress",
						Value: "in_progress",
					},
					{
						Label: "Completed",
						Value: "completed",
					},
				}).
				WrappClass("mt-4").
				Clearable(true).
				Required(true).
				ID("task-status-field"),

			NewCalendar().
				Label("Date").
				WrapClass("mt-4").
				Required(true).
				Value(&tf.data.Date).
				ID("task-date").
				OnUpdate(func(ctx app.Context, e app.Event) {
					tf.data.SetDate(tf.data.Date)
				}),

			NewCalendar().
				Label("Work time").
				WrapClass("mt-4").
				Required(true).
				Type(InputTypeDatetimeRange).
				Values(&tf.data.WorkBegin, &tf.data.WorkEnd).
				ID("work-time"),

			NewSelectField[int64]().
				Label("User").
				WrappClass("mt-4").
				Required(true).
				Value(&tf.data.UserID).
				Options(tf.userOptions()).
				ID("task-user"),
		)
}

func (tf *TaskForm) userOptions() (options []SelectOption[int64]) {
	for _, u := range tf.Ctx.Store.GetUsers() {
		options = append(options, SelectOption[int64]{
			Label: u.GetLabel(),
			Value: u.ID,
		})
	}

	return
}

func (tf *TaskForm) Render() app.UI {
	return NewDialog("taskDialog").
		Title("Task form").
		Visible(tf.visible).
		IsValid(tf.isValid()).
		OnOk(func(ctx app.Context, e app.Event) {
			if _, err := tf.Ctx.Store.SaveTask(tf.data); err != nil {
				tf.ShowNotification(ctx, "Error save user", err.Error())
			} else {
				tf.onHide(ctx)
			}
		}).
		OnDismiss(func(ctx app.Context, e app.Event) {
			ctx.NewAction("resetTask")
		}).
		Content(tf.El())
}

func (tf *TaskForm) isValid() bool {
	return tf.data.IsValid()
}

func (tf *TaskForm) onHide(ctx app.Context) {
	ctx.NewAction("hideModal")
	ctx.NewAction("resetTask")
}
