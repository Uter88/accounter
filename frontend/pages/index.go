package pages

import (
	"accounter/domain/user"
	"accounter/frontend/common"
	"accounter/frontend/components"
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
	return app.Form().Style("width", "500px").Class("d-flex flex-column").
		Body(
			app.Div().Class("d-flex flex-row align-items-center w-100").Body(
				app.Span().Class("material-symbols-outlined", "text-warning").Text("crown"),
				app.H5().Text("Registration").Class("px-1"),
			),

			// Name

			components.NewInputField[string]().
				Label("What is your name?").
				Value(i.form.Name).
				Rcv(&i.form.Name).
				WrapClass("mt-4").
				Clearable(true).
				Required(true).
				Autofocus(true).
				PrependIcon("timer_1").
				ID("name-field").
				OnInput(i.ValueTo(&i.form.Name)),

			components.NewInputField[string]().
				Label("What is your surname?").
				Value(i.form.Surname).
				Rcv(&i.form.Surname).
				Clearable(true).
				Required(true).
				PrependIcon("timer_2").
				ID("surname-field").
				OnInput(i.ValueTo(&i.form.Surname)),

			/* 			app.Label().
			   				For("name-field1").
			   				Text("* What is your name?").
			   				Class("mt-4", "text-secondary"),

			   			app.Div().
			   				Class("input-group").
			   				Body(
			   					app.Div().Class("input-group-prepend").Body(
			   						app.Span().Class("input-group-text", "material-symbols-outlined").ID("name").Text("timer_1"),
			   					),
			   					app.Input().
			   						Type("text").
			   						ID("name-field2").
			   						Class("form-control").
			   						Value(i.form.Name).
			   						Placeholder("").
			   						Attr("aria-label", "Username").
			   						Attr("aria-describedby", "name").
			   						AutoFocus(true).
			   						OnInput(i.ValueTo(&i.form.Name)),

			   					app.Div().Class("input-group-append").Body(
			   						app.If(len(i.form.Name) > 0, func() app.UI {
			   							return app.Span().Class("input-group-text", "material-symbols-outlined").Text("close").ID("close-name").
			   								OnClick(func(ctx app.Context, e app.Event) {
			   									i.form.ResetField("name")
			   								})
			   						}),
			   					),
			   				),

			   			app.Label().For("surname-field").Text("* What is your surname?").Class("mt-3", "text-secondary"),
			   			app.Div().Class("input-group").Body(
			   				app.Div().Class("input-group-prepend").Body(
			   					app.Span().Class("input-group-text", "material-symbols-outlined").ID("surname").Text("timer_2"),
			   				),
			   				app.Input().
			   					Type("text").
			   					ID("surname-field").
			   					Class("form-control").
			   					Value(i.form.Surname).
			   					Placeholder("").
			   					Attr("aria-label", "Surname").
			   					Attr("aria-describedby", "surname").
			   					OnInput(i.ValueTo(&i.form.Surname)),

			   				app.If(len(i.form.Surname) > 0, func() app.UI {
			   					return app.Span().Class("input-group-text", "material-symbols-outlined").Text("close").ID("close-surname").
			   						OnClick(func(ctx app.Context, e app.Event) {
			   							i.form.ResetField("surname")
			   						})
			   				}),
			   			),

			   			app.Label().For("patronymic-field").Text("What is your patronymic?").Class("mt-3", "text-secondary"),
			   			app.Div().Class("input-group").Body(
			   				app.Div().Class("input-group-prepend").Body(
			   					app.Span().Class("input-group-text", "material-symbols-outlined").ID("patronymic").Text("timer_3"),
			   				),
			   				app.Input().
			   					Type("text").
			   					ID("patronymic-field").
			   					Class("form-control").
			   					Value(i.form.Patronymic).
			   					Placeholder("").
			   					Attr("aria-label", "Patronymic").
			   					Attr("aria-describedby", "patronymic").
			   					OnInput(i.ValueTo(&i.form.Patronymic)),

			   				app.If(len(i.form.Patronymic) > 0, func() app.UI {
			   					return app.Span().Class("input-group-text", "material-symbols-outlined").Text("close").ID("close-patr").
			   						OnClick(func(ctx app.Context, e app.Event) {
			   							i.form.ResetField("patronymic")
			   						})
			   				}),
			   			),

			   			app.Label().For("login-field").Text("* Enter your login").Class("mt-3", "text-secondary"),
			   			app.Div().Class("input-group").Body(
			   				app.Div().Class("input-group-prepend").Body(
			   					app.Span().Class("input-group-text", "material-symbols-outlined").ID("login").Text("alternate_email"),
			   				),
			   				app.Input().
			   					Type("text").
			   					ID("login-field").
			   					Class("form-control").
			   					Value(i.form.Login).
			   					Placeholder("").
			   					Attr("aria-label", "Login").
			   					Attr("aria-describedby", "login").
			   					OnInput(i.ValueTo(&i.form.Login)),

			   				app.If(len(i.form.Login) > 0, func() app.UI {
			   					return app.Span().Class("input-group-text", "material-symbols-outlined").Text("close").ID("close-login").
			   						OnClick(func(ctx app.Context, e app.Event) {
			   							i.form.ResetField("login")
			   						})
			   				}),
			   			),

			   			app.Label().For("password-field").Text("* Enter your password").Class("mt-3", "text-secondary"),
			   			app.Div().Class("input-group").Body(
			   				app.Div().Class("input-group-prepend").Body(
			   					app.Span().Class("input-group-text", "material-symbols-outlined").ID("password").Text("password"),
			   				),
			   				app.Input().
			   					Type("password").
			   					ID("password-field").
			   					Class("form-control").
			   					Value(i.form.Password).
			   					Placeholder("").
			   					Attr("aria-label", "Password").
			   					Attr("aria-describedby", "password").
			   					OnInput(i.ValueTo(&i.form.Password)),

			   				app.If(len(i.form.Password) > 0, func() app.UI {
			   					return app.Span().Class("input-group-text", "material-symbols-outlined").Text("close").ID("close-passw").
			   						OnClick(func(ctx app.Context, e app.Event) {
			   							i.form.ResetField("password")
			   						})
			   				}),
			   			),

			   			app.Label().For("cost-field").Text("* What is your cost?").Class("mt-3", "text-secondary"),
			   			app.Div().Class("input-group").Body(
			   				app.Div().Class("input-group-prepend").Body(
			   					app.Span().Class("input-group-text", "material-symbols-outlined").ID("cost").Text("currency_ruble"),
			   				),
			   				app.Input().
			   					Type("number").
			   					ID("cost-field").
			   					Class("form-control").
			   					Value(i.form.PricePerHour).
			   					Placeholder("").
			   					Attr("aria-label", "Cost").
			   					Attr("aria-describedby", "cost").
			   					Min(1).
			   					OnInput(i.ValueTo(&i.form.PricePerHour)),
			   			), */

			app.Raw(`
				<a class="mt-3 icon-link icon-link-hover link-success link-underline-success link-underline-opacity-25" href="#">
				 I agree to the Terms of Service and Privacy Policy
					<svg xmlns="http://www.w3.org/2000/svg" class="bi" viewBox="0 0 16 16" aria-hidden="true">
						<path d="M1 8a.5.5 0 0 1 .5-.5h11.793l-3.147-3.146a.5.5 0 0 1 .708-.708l4 4a.5.5 0 0 1 0 .708l-4 4a.5.5 0 0 1-.708-.708L13.293 8.5H1.5A.5.5 0 0 1 1 8z"/>
					</svg>
				</a>
			`),

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
