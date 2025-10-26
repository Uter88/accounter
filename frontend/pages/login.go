package pages

import (
	"accounter/frontend/common"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type loginPage struct {
	common.BaseComponent
}

func NewLoginPage(ctx common.AppContext) *loginPage {
	return &loginPage{
		BaseComponent: common.NewBaseComponent(ctx),
	}
}

func (h *loginPage) Render() app.UI {
	btn := app.Button().
		Class("btn btn-primary").
		Style("width", "200px").
		Text("Request users").OnClick(func(ctx app.Context, e app.Event) {
		h.EnableNotifications(ctx, e)

		if err := h.Ctx.Store.RequestUsers(); err != nil {
			h.ShowNotification(ctx, "Error", err.Error())
		} else {
			h.ShowNotification(ctx, "Info", "Users loaded success!")
		}
	})

	rows := make([]app.UI, 0)

	for _, u := range h.Ctx.Store.GetUsers() {
		rows = append(rows, app.Tr().Body(
			app.Td().Text(u.Name),
			app.Td().Text(u.Surname),
			app.Td().Text(u.Patronymic),
		))
	}

	table := app.Table().
		Class("table table-striped table-hover table-bordered  w-100 mt-5").
		Body(
			app.THead().Body(
				app.Tr().Body(
					app.Th().Text("Name"),
					app.Th().Text("Surname"),
					app.Th().Text("Patronymic"),
				),
			),
			app.TBody().Body(rows...),
		)

	return app.Div().
		Class("container d-flex justify-content-center align-items-center").
		Body(
			app.Div().
				Class("card p-1 d-flex flex-column align-items-center").
				Style("width", "500px").
				Body(btn, table),
		)
}
