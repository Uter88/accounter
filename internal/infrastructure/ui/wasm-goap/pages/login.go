package pages

import (
	"accounter/internal/infrastructure/ui/wasm-goap/common"
	"accounter/internal/infrastructure/ui/wasm-goap/components"
	"accounter/internal/infrastructure/ui/wasm-goap/models"
	"accounter/pkg/utils"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type loginPage struct {
	app.Compo
	common.BaseComponent

	data models.LoginForm

	loaded bool
}

func NewLoginPage(ctx common.AppContext) *loginPage {
	return &loginPage{
		BaseComponent: common.NewBaseComponent(ctx),
		data:          models.NewLoginForm(),
	}
}

func (i *loginPage) OnMount(ctx app.Context) {
	if login, password, ok := i.Ctx.Store.LoadAuthData(ctx); ok {
		i.data.Login = login
		i.data.Password = password
	}

	i.loaded = true
}

func (i *loginPage) Render() app.UI {
	return app.If(i.loaded, i.El).Else(func() app.UI { return app.Span() })
}

func (i *loginPage) El() app.UI {
	return app.Div().
		Style("height", "100vh").
		Class("d-flex flex-row align-items-center w-100").
		Body(
			app.Div().
				Class("d-flex flex-row justify-content-start mx-auto align-self-center", "login-page").
				Body(
					app.Form().
						Style("width", "50%").
						Class("d-flex flex-column", "login-form").
						Body(
							app.Div().
								Class("d-flex flex-row align-items-center w-100").
								Body(
									app.Span().
										Class("material-symbols-outlined", "text-warning").
										Text(""),
									app.H4().
										Text("Sign In to Your Account").
										Class("px-1 text-bold"),
								),

							components.NewInputField[string]().
								Label("Enter your login").
								Value(&i.data.Login).
								WrapClass("mt-5").
								Clearable(true).
								Required(true).
								PrependIcon("alternate_email").
								Formater(utils.ClearEmail).
								Validator(func(ctx app.Context, s string) error {
									return utils.ValidEmail(s)
								}).
								ID("login-field"),

							components.NewInputField[string]().
								Label("Enter your password").
								Type("password").
								Value(&i.data.Password).
								WrapClass("mt-4").
								Clearable(true).
								Required(true).
								PrependIcon("password").
								ID("password-field"),

							components.NewCheckboxField[bool]().
								Label("Remember me").
								Value(&i.data.IsRemember).
								WrapClass("mt-4").
								ID("remember-field"),

							app.Button().
								Text("Sign In").
								Type("button").
								Class("mt-5 btn btn-primary btn-lg").
								Disabled(!i.data.Validate(true)).
								OnClick(func(ctx app.Context, e app.Event) {
									if err := i.Ctx.Store.LoginByCredentials(ctx, i.data); err != nil {
										i.ShowNotification(ctx, "Error", err.Error())

									} else {
										i.data.Reset()
										ctx.Navigate("/index")
									}
								}),

							app.Raw(`
								<a class="mt-3 icon-link icon-link-hover link-secondary link-underline-light link-underline-opacity-25" href="/registration">
								I have already an account: <b>Sign Up</b>
									<svg xmlns="http://www.w3.org/2000/svg" class="bi" viewBox="0 0 16 16" aria-hidden="true">
										<path d="M1 8a.5.5 0 0 1 .5-.5h11.793l-3.147-3.146a.5.5 0 0 1 .708-.708l4 4a.5.5 0 0 1 0 .708l-4 4a.5.5 0 0 1-.708-.708L13.293 8.5H1.5A.5.5 0 0 1 1 8z"/>
									</svg>
								</a>
							`),
						)))

}
