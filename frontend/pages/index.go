package pages

import (
	"accounter/domain/task"
	"accounter/domain/user"
	"accounter/frontend/common"
	"accounter/frontend/components"
	"accounter/frontend/models"
	"math/rand/v2"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type indexPage struct {
	app.Compo
	common.BaseComponent
}

func NewIndexPage(ctx common.AppContext) *indexPage {
	return &indexPage{
		BaseComponent: common.NewBaseComponent(ctx),
	}
}

func (inp *indexPage) requestUsers(ctx app.Context) {
	inp.EnableNotifications(ctx)

	if err := inp.Ctx.Store.RequestUsers(); err != nil {
		inp.ShowNotification(ctx, "Error", err.Error())
	} else {
		inp.ShowNotification(ctx, "Info", "Users loaded success!")
	}
}

func (inp *indexPage) requestTasks(ctx app.Context) {
	inp.EnableNotifications(ctx)

	if err := inp.Ctx.Store.RequestTasks(); err != nil {
		inp.ShowNotification(ctx, "Error", err.Error())
	} else {
		inp.ShowNotification(ctx, "Info", "Tasks loaded success!")
	}
}

func (inp *indexPage) OnMount(ctx app.Context) {
	if !inp.Ctx.Store.CheckAuth(ctx) {
		ctx.Navigate("/login")
		return
	}

	inp.requestUsers(ctx)
	inp.requestTasks(ctx)
}

func (inp *indexPage) GroupBtn() app.HTMLDiv {
	btnIcon := components.NewBtnIcon("logout").
		Text("Exit").
		Tooltip("Exit from system").
		OnClick(func(ctx app.Context, e app.Event) {
			inp.Ctx.Store.Logout(ctx)
		})

	return app.Div().
		Body(
			btnIcon,
		)
}

func (inp *indexPage) initChart(opts *models.ChartOptions) *components.Chart {
	chart := components.NewChart("bar-chart").
		Width("500px").
		Height("400px").
		SetOptions(opts).
		Lazy(true)

	return chart
}

func getChartOptions() *models.ChartOptions {
	var data = models.Array{5, 20, 36, 10, 10, 20}

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	return models.NewChartOptions().
		WithTitle("Bar chart").
		WithLegend("sales").
		WithXAxis(models.XAxis{Data: models.Array{"Shirts", "Cardigans", "Chiffons", "Pants", "Heels", "Socks"}}).
		WithSeries(models.Series{Name: "sales", Type: "bar", Data: data})
}

func (inp *indexPage) Render() app.UI {
	userForm := components.NewUserForm(inp.Ctx, false)
	taskForm := components.NewTaskForm(inp.Ctx)

	usersTable := components.NewUserList(inp.Ctx, inp.Ctx.Store.GetUsers()).
		Loading(inp.Ctx.Store.GetUsersLoading()).
		OnRequest(inp.requestUsers).
		OnEdit(func(ctx app.Context, u user.User) {
			ctx.NewActionWithValue("setUser", u)
		}).
		OnAdd(func(ctx app.Context) {
			ctx.NewActionWithValue("setUser", user.User{})
		})

	tasksTable := components.NewTaskList(inp.Ctx, inp.Ctx.Store.GetTasks()).
		OnRequest(func(ctx app.Context, p *task.TaskParams) {
			inp.Ctx.Store.SetTaskParams(p)
			inp.requestTasks(ctx)
		}).
		OnEdit(func(ctx app.Context, t task.Task) {
			ctx.NewActionWithValue("setTask", t)
		}).
		OnAdd(func(ctx app.Context) {
			t := task.NewTask()
			params := inp.Ctx.Store.GetTaskParams()
			tasks := inp.Ctx.Store.GetTasks()

			t.SetDate(params.DateStart)

			if len(params.Users) > 0 {
				t.UserID = params.Users[0]

				if l := len(tasks); l > 0 && tasks[l-1].UserID == t.UserID {
					t.SetDate(tasks[l-1].Date)
				}
			}

			ctx.NewActionWithValue("setTask", t)
		}).
		OnDelete(func(ctx app.Context, t task.Task) {
			if err := inp.Ctx.Store.RemoveTask(t); err != nil {
				inp.ShowNotification(ctx, "Error", err.Error())
			}
		})

	chart := inp.initChart(getChartOptions())

	return app.Main().
		Class("d-flex w-100 h-100 flex-column").
		Body(
			// Header
			app.Header().
				Class("d-flex flex-row w-100 py-2 px-3 justify-content-end").
				Body(inp.GroupBtn()),

			// Body
			app.Div().
				Class("d-flex flex-row w-100 p-3").
				Body(
					app.Div().
						Class("d-flex h-100 col-5 flex-column").
						Body(
							app.Div().
								Class("card p-1 d-flex h-50 flex-column align-items-center mx-3").
								Style("border", "1px solid red").
								Body(userForm, usersTable),

							app.Div().
								Class("card p-1 d-flex h-50 flex-column align-items-center mt-3 mx-3").
								Body(
									components.NewBtnIcon("bar_chart").
										OnClick(func(ctx app.Context, e app.Event) {
											chart.SetOptions(getChartOptions()).Draw(ctx)
										}),

									chart,
								),
						),

					app.Div().
						Class("card p-1 d-flex flex-column align-items-center h-100 mx-1").
						Style("border", "1px solid red").
						Body(taskForm, tasksTable),
				),
		)
}
