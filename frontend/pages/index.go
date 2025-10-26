package pages

import (
	"accounter/domain/user"
	"accounter/frontend/common"
	"log"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type indexPage struct {
	common.BaseComponent

	form user.User
}

func NewIndexPage(ctx common.AppContext) *indexPage {
	return &indexPage{
		BaseComponent: common.NewBaseComponent(ctx),
	}
}

func (i *indexPage) Render() app.UI {
	return app.Form().Style("width", "500px").Class("d-flex flex-column", "px-5").
		Body(
			app.H5().Text("Registration"),
			app.Input().
				Type("text").
				ID("login-field").
				Class("mt-3 form-control").
				Value(i.form.Name).
				Placeholder("What is your name?").
				AutoFocus(true).
				OnInput(i.ValueTo(&i.form.Name)),

			app.Input().
				Type("text").
				Class("mt-3 form-control").
				Value(i.form.Surname).
				Placeholder("What is your surname?").
				OnInput(i.ValueTo(&i.form.Surname)),

			app.Input().
				Type("text").
				Class("mt-3 form-control").
				Value(i.form.Patronymic).
				Placeholder("What is your patronymic?").
				OnInput(i.ValueTo(&i.form.Patronymic)),

			app.Input().
				Type("login").
				Class("mt-3 form-control").
				Value(i.form.Login).
				Placeholder("What is your login?").
				OnInput(i.ValueTo(&i.form.Login)),

			app.Input().
				Type("password").
				Class("mt-3 form-control").
				Value(i.form.Password).
				Placeholder("What is your password?").
				OnInput(i.ValueTo(&i.form.Password)),

			app.Input().
				Type("number").
				Class("mt-3 form-control").
				Min(1).
				Value(i.form.PricePerHour).
				Placeholder("What is your cost?").
				OnInput(i.ValueTo(&i.form.PricePerHour)),

			app.Button().
				Text("Save").
				Class("mt-3 btn btn-primary").
				Disabled(!i.form.IsValid()).
				OnClick(func(ctx app.Context, e app.Event) {
					if err := i.Ctx.Store.SaveUser(i.form); err != nil {
						log.Println(err.Error())
					}
					i.form.Reset()
				}),
		)
}
