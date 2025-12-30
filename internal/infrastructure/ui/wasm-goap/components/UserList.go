package components

import (
	"accounter/internal/domain/user"
	"accounter/internal/infrastructure/ui/wasm-goap/common"
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserList struct {
	app.Compo
	common.BaseComponent

	Users    user.Users
	selected []int64
}

func NewUserList(ctx common.AppContext, users user.Users) *UserList {
	return &UserList{
		BaseComponent: common.NewBaseComponent(ctx),
		Users:         users,
	}
}

func (ul *UserList) OnMount(ctx app.Context) {
	ul.onRequest(ctx)
}

func (ul *UserList) onRequest(ctx app.Context) {
	ul.EnableNotifications(ctx)

	if err := ul.Ctx.Store.RequestUsers(); err != nil {
		ul.ShowNotification(ctx, "Error", err.Error())
	}
}

func (ul *UserList) onEdit(ctx app.Context, u *user.User) {
	if u == nil {
		u = &user.User{}
	}

	ctx.NewActionWithValue("setUser", u)
}

func (ul *UserList) isLoading() bool {
	return ul.Ctx.Store.GetTasksLoading()
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
							ul.onEdit(ctx, &ul.Users[i])
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
					ul.onEdit(ctx, nil)
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
		NewLoading(ul.isLoading()),
		toolbar,
		table,
	)
}
