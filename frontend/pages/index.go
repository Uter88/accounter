package pages

import (
	"accounter/domain/user"
	"accounter/frontend/common"
	"accounter/frontend/components"

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

func (inp *indexPage) OnMount(ctx app.Context) {
	if !inp.Ctx.Store.CheckAuth(ctx) {
		ctx.Navigate("/login")
		return
	}

	inp.requestUsers(ctx)
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

func (inp *indexPage) Render() app.UI {
	usersTable := components.NewUserList(inp.Ctx, inp.Ctx.Store.GetUsers()).
		Loading(inp.Ctx.Store.GetUsersLoading()).
		OnRequest(inp.requestUsers).
		OnEdit(func(ctx app.Context, u user.User) {
			ctx.NewActionWithValue("setUser", u)
		}).
		OnAdd(func(ctx app.Context) {
			ctx.NewActionWithValue("setUser", user.User{})
		})

	return app.Main().
		Class("d-flex w-100 h-100 flex-column").
		Body(
			// Header
			app.Header().
				Class("d-flex flex-row w-100 py-2 px-3 justify-content-end").
				Body(inp.GroupBtn()),

			// Body
			app.Div().
				Class("d-flex flex-row h-100 w-100 p-3").
				Body(
					app.Div().
						Class("d-flex h-100 col-5 flex-column").
						Body(
							app.Div().
								Class("card p-1 d-flex h-50 flex-column align-items-center mx-3").
								Style("border", "1px solid red").
								Body(
									components.NewUserForm(inp.Ctx, false),
									usersTable,
								),

							components.NewTasksChart(inp.Ctx.Store.GetTasks(), inp.Ctx.Store.GetTaskParams()),
						),

					app.Div().
						Class("card p-1 d-flex flex-column align-items-center h-100 mx-1").
						Style("border", "1px solid red").
						Body(
							components.NewTaskForm(inp.Ctx),
							components.NewTaskList(inp.Ctx, inp.Ctx.Store.GetTasks()),
						),
				),
		)
}
