package components

import (
	"accounter/internal/domain/task"
	"accounter/internal/infrastructure/ui/wasm-goap/common"
	"accounter/pkg/tools"
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type TaskList struct {
	app.Compo
	common.BaseComponent

	Tasks  task.Tasks
	params *task.TaskParams
}

func NewTaskList(ctx common.AppContext, tasks task.Tasks) *TaskList {
	return &TaskList{
		BaseComponent: common.NewBaseComponent(ctx),
		Tasks:         tasks,
		params:        ctx.Store.GetTaskParams(),
	}
}

func (tl *TaskList) OnMount(ctx app.Context) {
	tl.onRequest(ctx)
}

func (tl *TaskList) isLoading() bool {
	return tl.Ctx.Store.GetTasksLoading()
}

func (tl *TaskList) Render() app.UI {
	rows := make([]app.UI, len(tl.Tasks))

	for i, t := range tl.Tasks {
		rows[i] = app.Tr().Body(
			app.Td().Body(
				app.Div().Class("d-flex flex-row justify-content-center").Body(
					NewBtnIcon("edit").
						Tooltip("Edit task").
						Target("#taskDialog").
						OnClick(func(ctx app.Context, e app.Event) {
							tl.onEdit(ctx, &tl.Tasks[i])
						}),

					NewBtnIcon("delete").
						Tooltip("Delete").
						BtnClass("mx-1").
						OnClick(func(ctx app.Context, e app.Event) {
							tl.onDelete(ctx, t)
						}),
				),
			),
			app.Td().Text(t.FormatDate()),
			app.Td().Text(t.TaskID),
			app.Td().Text(t.Description),
			app.Td().Text(t.UserLabel).Class("wrap-cell"),
			app.Td().Text(t.Status),
			app.Td().Text(t.FormatWorkBegin()),
			app.Td().Text(t.FormatWorkEnd()),
			app.Td().Text(t.FormatDuration(false)),
			app.Td().Text(t.FormatPrice()),
		)
	}

	toolbar := app.Div().
		Class("p-3 d-flex flex-row align-items-center").
		Body(
			NewCalendar().
				Values(&tl.params.DateStart, &tl.params.DateEnd).
				OnUpdate(func(ctx app.Context, e app.Event) {
					tl.onRequest(ctx)
				}),
			NewBtnIcon("refresh").
				Tooltip("refresh").
				OnClick(func(ctx app.Context, e app.Event) {
					tl.onRequest(ctx)
				}),
			NewBtnIcon("add").
				Tooltip("add").
				Target("#taskDialog").
				OnClick(func(ctx app.Context, e app.Event) {
					tl.onEdit(ctx, nil)
				}),

			NewBtnIcon("download").
				Tooltip("export").
				Disabled(tl.Tasks.Empty()).
				OnClick(func(ctx app.Context, e app.Event) {
					tl.onExport(ctx, tools.FileFormatHTML)
				}),
		)
	table := app.Span().
		Class("table-wrapper").
		Body(
			app.Table().
				Class("table table-striped table-hover table-bordered w-100 sticky-header sticky-footer").
				Style("vertical-align", "middle").
				Body(
					app.THead().Body(
						app.Tr().Body(
							app.Th().Text(""),
							app.Th().Text("Date"),
							app.Th().Text("Task ID"),
							app.Th().Text("Description"),
							app.Th().Text("User"),
							app.Th().Text("Status"),
							app.Th().Text("Work begin"),
							app.Th().Text("Work end"),
							app.Th().Text("Duration"),
							app.Th().Text("Cost"),
						),
					),
					app.TBody().Body(rows...),
					app.TFoot().Body(tl.getSummary()),
				),
		)

	return app.Div().
		Class("w-100 h-100").
		Style("overflow", "hidden").
		Body(
			NewLoading(tl.isLoading()),
			toolbar,
			table,
		)
}

func (tl *TaskList) onRequest(ctx app.Context) {
	tl.Ctx.Store.SetTaskParams(tl.params)

	if err := tl.Ctx.Store.RequestTasks(ctx); err != nil {
		tl.ShowNotification(ctx, "Error", err.Error())
	}
}

func (tl *TaskList) onEdit(ctx app.Context, t *task.Task) {
	if t == nil {
		nt := task.NewTask()
		t = &nt
		params := tl.Ctx.Store.GetTaskParams()
		tasks := tl.Ctx.Store.GetTasks()

		t.SetDate(params.DateStart)

		if len(params.Users) > 0 {
			t.UserID = params.Users[0]

			if l := len(tasks); l > 0 && tasks[l-1].UserID == t.UserID {
				t.SetDate(tasks[l-1].Date)
			}
		}
	}

	ctx.NewActionWithValue("setTask", *t)
}

func (tl *TaskList) onDelete(ctx app.Context, t task.Task) {
	if !tl.ShowConfirm(ctx, fmt.Sprintf("Delete task %s?", t.Description)) {
		return
	}

	if err := tl.Ctx.Store.RemoveTask(ctx, t); err != nil {
		tl.ShowNotification(ctx, "Error", err.Error())
	}
}

func (tl *TaskList) onExport(ctx app.Context, format tools.FileFormat) {
	api := tl.Ctx.Store.ExportTasks(format)
	ctx.Navigate(api)
}

func (tl *TaskList) getSummary() app.UI {
	var (
		dur   time.Duration
		price float32
		cnt   int
	)

	for _, t := range tl.Tasks {
		dur += t.GetDuration()
		price += t.GetPrice()
		cnt++
	}

	return app.Tr().Body(
		app.Td().Text(cnt).Attr("align", "center").Class("text-bold"),
		app.Td(), app.Td(), app.Td(), app.Td(), app.Td(), app.Td(),
		app.Td(),
		app.Td().Text(tools.FormatDuration(dur, false)).Class("text-bold"),
		app.Td().Text(fmt.Sprintf("%.2f", price)).Class("text-bold"),
	)
}
