package components

import (
	"accounter/domain/user"
	"accounter/frontend/common"
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserList struct {
	app.Compo
	common.BaseComponent

	onRequest func(ctx app.Context)
	onAdd     func(ctx app.Context)
	onEdit    func(ctx app.Context, u user.User)
	Users     user.Users
	loading   bool
	selected  []int64
}

func NewUserList(ctx common.AppContext, users user.Users) *UserList {
	return &UserList{
		BaseComponent: common.NewBaseComponent(ctx),
		Users:         users,
		onAdd:         func(ctx app.Context) {},
		onEdit:        func(ctx app.Context, u user.User) {},
		onRequest:     func(ctx app.Context) {},
	}
}

func (ul *UserList) OnRequest(cb func(ctx app.Context)) *UserList {
	ul.onRequest = cb
	return ul
}

func (ul *UserList) OnAdd(cb func(ctx app.Context)) *UserList {
	ul.onAdd = cb
	return ul
}

func (ul *UserList) OnEdit(cb func(ctx app.Context, u user.User)) *UserList {
	ul.onEdit = cb
	return ul
}

func (ul *UserList) Loading(v bool) *UserList {
	ul.loading = v
	return ul
}

func (ul *UserList) Render() app.UI {
	rows := make([]app.UI, len(ul.Users))

	for i, u := range ul.Users {
		rows[i] = app.Tr().Body(
			app.Td().Body(
				app.Div().Class("d-flex flex-row justify-content-center").Body(
					NewCheckboxField[int64]().
						Value(&u.ID).
						Values(&ul.selected).
						OnUpdate(func(ctx app.Context, e app.Event) {
							params := ul.Ctx.Store.GetTaskParams()
							params.Users = ul.selected
							ul.Ctx.Store.SetTaskParams(params)
							ul.Ctx.Store.RequestTasks(ctx)
						}),

					NewBtnIcon("edit").
						Tooltip("Edit user").
						Target("#userDialog").
						OnClick(func(ctx app.Context, e app.Event) {
							ul.onEdit(ctx, ul.Users[i])
						}),
					NewBtnIcon("delete").
						Tooltip("Delete").
						BtnClass("mx-1").
						OnClick(func(ctx app.Context, e app.Event) {
							fmt.Println("on delete ", u.ID)
						}),
				),
			),
			app.Td().Text(u.Login),
			app.Td().Text(u.Name),
			app.Td().Text(u.Surname),
			app.Td().Text(u.Patronymic),
			app.Td().Text(u.PricePerHour),
		)
	}

	toolbar := app.Div().
		Class("p-3 d-flex flex-row").
		Body(
			NewBtnIcon("refresh").
				Tooltip("refresh").
				OnClick(func(ctx app.Context, e app.Event) {
					ul.onRequest(ctx)
				}),
			NewBtnIcon("add").
				Tooltip("add").
				Target("#userDialog").
				OnClick(func(ctx app.Context, e app.Event) {
					ul.onAdd(ctx)
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
					app.Th().Text("Login"),
					app.Th().Text("Name"),
					app.Th().Text("Surname"),
					app.Th().Text("Patronymic"),
					app.Th().Text("Price"),
				),
			),
			app.TBody().Body(rows...),
		)

	return app.Div().Body(
		toolbar,
		table,
	)
}
