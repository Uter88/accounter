package pages

import (
	"accounter/domain/user"
	"accounter/frontend/common"
	"accounter/frontend/components"
	"log"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type registrationPage struct {
	common.BaseComponent

	form user.User
}

func NewRegistrationPage(ctx common.AppContext) *registrationPage {
	return &registrationPage{
		BaseComponent: common.NewBaseComponent(ctx),
	}
}

func (i *registrationPage) Render() app.UI {
	return app.Div().
		Style("height", "100vh").
		Class("d-flex flex-row align-items-center w-100").
		Body(
			app.Div().
				Class("d-flex flex-row justify-content-start mx-auto align-self-center", "login-page").Body(
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
									Text("Create account").
									Class("px-1 text-bold"),
							),

						components.NewInputField[string]().
							Label("What is your name?").
							Value(&i.form.Name).
							WrapClass("mt-4").
							Clearable(true).
							Required(true).
							Autofocus(true).
							PrependIcon("timer_1").
							ID("name-field"),

						components.NewInputField[string]().
							Label("What is your surname?").
							Value(&i.form.Surname).
							WrapClass("mt-3").
							Clearable(true).
							Required(true).
							PrependIcon("timer_2").
							ID("surname-field"),

						components.NewInputField[string]().
							Label("What is your patronymic?").
							Value(&i.form.Patronymic).
							WrapClass("mt-3").
							Clearable(true).
							Required(false).
							PrependIcon("timer_3").
							ID("patronymic-field"),

						components.NewInputField[string]().
							Label("Enter your login").
							Value(&i.form.Login).
							WrapClass("mt-3").
							Clearable(true).
							Required(true).
							PrependIcon("alternate_email").
							ID("login-field"),

						components.NewInputField[string]().
							Label("Enter your password").
							Type("password").
							Value(&i.form.Password).
							WrapClass("mt-3").
							Clearable(true).
							Required(true).
							PrependIcon("password").
							ID("password-field"),

						components.NewInputField[float32]().
							Label("Enter price of your job").
							Type("number").
							WrapClass("mt-3").
							Value(&i.form.PricePerHour).
							PrependIcon("currency_ruble").
							Min(1).
							Step(0.01).
							ID("cost-field"),

						components.NewCheckboxField().
							Label("Remember me").
							Checked(&i.form.IsRemember).
							WrapClass("mt-3").
							ID("remember-field"),

						components.NewCheckboxField().
							Label("I agree to the Terms of Service and Privacy Policy").
							Checked(&i.form.IsAccept).
							WrapClass("mt-1").
							Required(true).
							ID("accept-field"),

						app.Button().
							Text("Save").
							Class("mt-3 btn btn-primary btn-lg").
							Disabled(!i.form.IsValid(false)).
							OnClick(func(ctx app.Context, e app.Event) {
								if err := i.Ctx.Store.SaveUser(i.form); err != nil {
									log.Println(err.Error())
								}
								i.form.Reset()
							}),

						app.Raw(`
							<a class="mt-3 icon-link icon-link-hover link-secondary link-underline-light link-underline-opacity-25" href="/login">
							I have already an account: <b>Sign In</b>
								<svg xmlns="http://www.w3.org/2000/svg" class="bi" viewBox="0 0 16 16" aria-hidden="true">
									<path d="M1 8a.5.5 0 0 1 .5-.5h11.793l-3.147-3.146a.5.5 0 0 1 .708-.708l4 4a.5.5 0 0 1 0 .708l-4 4a.5.5 0 0 1-.708-.708L13.293 8.5H1.5A.5.5 0 0 1 1 8z"/>
								</svg>
							</a>
					`),
					)))
}
