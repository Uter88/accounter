package pages

import (
	"accounter/frontend/common"
	"accounter/frontend/components"
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type listPage struct {
	common.BaseComponent
}

func NewListPage(ctx common.AppContext) *listPage {
	return &listPage{
		BaseComponent: common.NewBaseComponent(ctx),
	}
}

func (h *listPage) Render() app.UI {
	btnReq := app.Button().
		Class("btn btn-outline-primary px-1").
		Style("width", "auto").
		Text("Request users").OnClick(func(ctx app.Context, e app.Event) {
		h.EnableNotifications(ctx, e)

		if err := h.Ctx.Store.RequestUsers(); err != nil {
			h.ShowNotification(ctx, "Error", err.Error())
		} else {
			h.ShowNotification(ctx, "Info", "Users loaded success!")
		}
	})

	exitIcon := components.NewIcon("logout").Class("text-secondary")

	exitBtn := app.Div().
		Class("d-flex flex-row align-items-center px-1").
		Body(
			exitIcon.Class("px-1"),
			app.Button().
				Class("btn btn-outline-light text-secondary btn-flat").
				Text("Exit"),
		).
		Role("button").
		OnClick(func(ctx app.Context, e app.Event) {
			ctx.Navigate("/login")
		})

	btnGroup := app.Div().
		Class("d-flex flex-row align-items-center w-100 justify-content-between").
		Body(btnReq, exitBtn)

	rows := make([]app.UI, 0)

	for _, u := range h.Ctx.Store.GetUsers() {
		rows = append(rows, app.Tr().Body(
			app.Td().Body(
				components.NewIcon("edit").
					El().
					Role("button").
					OnClick(func(ctx app.Context, e app.Event) {
						fmt.Println("on edit ", u.Surname)
					}),

				components.NewIcon("eye_tracking").
					El().
					Role("button").
					OnClick(func(ctx app.Context, e app.Event) {
						fmt.Println("on preview ", u.Surname)
					}),
			),
			app.Td().Text(u.Login),
			app.Td().Text(u.Name),
			app.Td().Text(u.Surname),
			app.Td().Text(u.Patronymic),
			app.Td().Text(u.PricePerHour),
		))
	}

	table := app.Table().
		Class("table table-striped table-hover table-bordered  w-100 mt-5").
		Body(
			app.THead().Body(
				app.Tr().Body(
					app.Th().Text(""),
					app.Th().Text("Login"),
					app.Th().Text("Name"),
					app.Th().Text("Surname"),
					app.Th().Text("Patronymic"),
					app.Th().Text("Price"),
				),
			),
			app.TBody().Body(rows...),
		)

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
