package pages

import (
	"accounter/internal/infrastructure/ui/wasm-goap/common"
	"accounter/internal/infrastructure/ui/wasm-goap/components"

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

func (inp *indexPage) OnMount(ctx app.Context) {
	if !inp.Ctx.Store.CheckAuth(ctx) {
		ctx.Navigate("/login")
		return
	}

	inp.Ctx.Websocket.Connect(inp.Ctx.Store.GetUser().Tokens.AccessToken)
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
	return app.Main().
		Class("d-flex w-100 h-100 flex-column").
		Body(
			// Header
			app.Header().
				Class("d-flex flex-row w-100 py-2 px-3 justify-content-between").
				Body(
					app.H5().Text("Personal Area").Class("mx-3"),
					inp.GroupBtn()),

			// Body
			app.Div().
				Class("d-flex flex-row h-100 w-100 p-3").
				Body(
					app.Div().
						Class("d-flex h-100 col-5 flex-column").
						Body(
							app.Div().
								Class("card p-1 d-flex h-50 flex-column align-items-center mx-3").
								Body(
									components.NewUserForm(inp.Ctx, false),
									components.NewUserList(inp.Ctx, inp.Ctx.Store.GetUsers()),
								),

							components.NewTasksChart(inp.Ctx.Store.GetTasks(), inp.Ctx.Store.GetTaskParams()),
						),

					app.Div().
						Class("card p-1 d-flex flex-column align-items-center h-100 mx-1").
						Body(
							components.NewTaskForm(inp.Ctx),
							components.NewTaskList(inp.Ctx, inp.Ctx.Store.GetTasks()),
						),
				),
		)
}
