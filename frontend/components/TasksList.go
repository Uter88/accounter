package components

import (
	"accounter/domain/task"
	"accounter/frontend/common"
	"accounter/pkg/tools"
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type TaskList struct {
	app.Compo
	common.BaseComponent

	onRequest func(ctx app.Context, p *task.TaskParams)
	onAdd     func(ctx app.Context)
	onEdit    func(ctx app.Context, t task.Task)
	onDelete  func(ctx app.Context, t task.Task)
	onPreview func(ctx app.Context, t task.Task)
	Tasks     task.Tasks

	params *task.TaskParams
}

func NewTaskList(ctx common.AppContext, tasks task.Tasks) *TaskList {
	return &TaskList{
		Tasks:     tasks,
		onAdd:     func(ctx app.Context) {},
		onEdit:    func(ctx app.Context, t task.Task) {},
		onDelete:  func(ctx app.Context, t task.Task) {},
		onPreview: func(ctx app.Context, t task.Task) {},
		onRequest: func(ctx app.Context, p *task.TaskParams) {},
		params:    ctx.Store.GetTaskParams(),
	}
}

func (tl *TaskList) OnRequest(cb func(ctx app.Context, p *task.TaskParams)) *TaskList {
	tl.onRequest = cb
	return tl
}

func (tl *TaskList) OnAdd(cb func(ctx app.Context)) *TaskList {
	tl.onAdd = cb
	return tl
}

func (tl *TaskList) OnEdit(cb func(ctx app.Context, t task.Task)) *TaskList {
	tl.onEdit = cb
	return tl
}

func (tl *TaskList) OnDelete(cb func(ctx app.Context, t task.Task)) *TaskList {
	tl.onDelete = cb
	return tl
}

func (tl *TaskList) OnPreview(cb func(ctx app.Context, t task.Task)) *TaskList {
	tl.onPreview = cb
	return tl
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
							tl.onEdit(ctx, tl.Tasks[i])
						}),

					NewBtnIcon("delete").
						Tooltip("Delete").
						BtnClass("mx-1").
						OnClick(func(ctx app.Context, e app.Event) {
							tl.onDelete(ctx, tl.Tasks[i])
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
					tl.onRequest(ctx, tl.params)
				}),
			NewBtnIcon("refresh").
				Tooltip("refresh").
				OnClick(func(ctx app.Context, e app.Event) {
					tl.onRequest(ctx, tl.params)
				}),
			NewBtnIcon("add").
				Tooltip("add").
				Target("#taskDialog").
				OnClick(func(ctx app.Context, e app.Event) {
					tl.onAdd(ctx)
				}),
		)

	table := app.Table().
		Class("table table-striped table-hover table-bordered w-100").
		Style("table-layout", "fixed").
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
		)

	return app.Div().
		Class("w-100 h-100").
		Style("overflow", "auto").
		Style("max-height", "800px").
		Body(
			toolbar,
			table,
		)
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
