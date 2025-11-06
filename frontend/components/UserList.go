package components

import (
	"accounter/domain/user"
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserList struct {
	app.Compo

	Users user.Users
}

func NewUserList(users user.Users) *UserList {
	return &UserList{Users: users}
}

func (ul *UserList) Render() app.UI {
	rows := make([]app.UI, 0)

	for _, u := range ul.Users {
		rows = append(rows, app.Tr().Body(
			app.Td().Body(
				app.Div().Class("d-flex flex-row align-items-center").Body(

					NewIcon("edit", "Edit user").
						AlignClass("tooltip-top").
						El().
						Role("button").
						OnClick(func(ctx app.Context, e app.Event) {
							fmt.Println("on edit ", u.Surname)
						}),

					NewIcon("eye_tracking", "Preview").
						Class("mx-1").
						El().
						Role("button").
						OnClick(func(ctx app.Context, e app.Event) {
							fmt.Println("on preview ", u.Surname)
						}),
				),
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

	return table
}
