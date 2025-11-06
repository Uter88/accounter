package pages

import (
	"accounter/frontend/common"
	"accounter/frontend/components"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type indexPage struct {
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
	btnReq := app.Button().
		Class("btn btn-outline-primary px-1").
		Style("width", "auto").
		Text("reload").
		OnClick(func(ctx app.Context, e app.Event) {
			inp.requestUsers(ctx)
		})

	btnIcon := components.NewBtnIcon().
		OnClick(func(ctx app.Context, e app.Event) {
			inp.Ctx.Store.Logout(ctx)
		})

	btnGroup := app.Div().
		Class("d-flex flex-row align-items-center w-100 justify-content-between").
		Body(btnReq, btnIcon)

	return btnGroup
}

func (inp *indexPage) Render() app.UI {
	btnGroup := inp.GroupBtn()
	table := components.NewUserList(inp.Ctx.Store.GetUsers())

	return app.Div().
		Class("d-flex w-100 flex-row").
		Body(
			app.Div().
				Class("card p-1 d-flex flex-column align-items-center w-50 mx-3").
				Style("border", "1px solid red").
				Body(btnGroup, table),

			app.Div().
				Class("card p-1 d-flex flex-column align-items-center w-50 mx-3").
				Style("border", "1px solid red").
				Body(),
		)
}
