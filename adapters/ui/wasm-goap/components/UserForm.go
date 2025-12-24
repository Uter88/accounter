package components

import (
	"accounter/adapters/ui/wasm-goap/common"
	"accounter/adapters/ui/wasm-goap/models"

	"accounter/internal/domain/user"
	"accounter/pkg/tools"
	"errors"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type UserForm struct {
	app.Compo
	common.BaseComponent

	isRegister bool
	formClass  string
	data       models.LoginForm
	visible    bool
}

func NewUserForm(ctx common.AppContext, register bool) *UserForm {
	return &UserForm{
		isRegister:    register,
		BaseComponent: common.NewBaseComponent(ctx),
		data:          models.NewLoginForm(),
	}
}

func (uf *UserForm) OnMount(ctx app.Context) {
	ctx.Handle("setUser", func(ctx app.Context, a app.Action) {
		user := a.Value.(user.User)
		uf.data.User = &user
		uf.Open()
	})

	ctx.Handle("resetUser", func(ctx app.Context, a app.Action) {
		uf.data.User = &user.User{}
		uf.Hide()
	})
}

func (uf *UserForm) FormClass(c string) *UserForm {
	uf.formClass = c
	return uf
}

func (uf *UserForm) Open() *UserForm {
	uf.visible = true
	return uf
}

func (uf *UserForm) Hide() *UserForm {
	uf.visible = false
	return uf
}

func (uf *UserForm) El() app.UI {
	return app.Form().
		Class("d-flex flex-column", uf.formClass).
		Body(
			NewInputField[string]().
				Label("What is your name?").
				Value(&uf.data.Name).
				WrapClass("mt-4").
				Clearable(true).
				Required(true).
				Autofocus(true).
				PrependIcon("timer_1").
				ID("name-field"),

			NewInputField[string]().
				Label("What is your surname?").
				Value(&uf.data.Surname).
				WrapClass("mt-3").
				Clearable(true).
				Required(true).
				PrependIcon("timer_2").
				ID("surname-field"),

			NewInputField[string]().
				Label("What is your patronymic?").
				Value(&uf.data.Patronymic).
				WrapClass("mt-3").
				Clearable(true).
				Required(false).
				PrependIcon("timer_3").
				ID("patronymic-field"),

			NewInputField[string]().
				Label("Enter your login").
				Value(&uf.data.Login).
				WrapClass("mt-3").
				Clearable(true).
				Required(true).
				PrependIcon("alternate_email").
				Validator(func(ctx app.Context, s string) error {
					if err := tools.ValidEmail(s); err != nil {
						return err
					}

					if ok, err := uf.Ctx.Store.CheckUniqueLogin(s, uf.data.ID); err != nil {
						uf.ShowNotification(ctx, "Error", err.Error())
						return err
					} else if ok {
						return errors.New("user with this login is exists")
					}

					return nil
				}).
				ID("login-field"),

			NewInputField[string]().
				Label("Enter your password").
				Type("password").
				Value(&uf.data.Password).
				WrapClass("mt-3").
				Clearable(true).
				Required(true).
				PrependIcon("password").
				ID("password-field"),

			NewInputField[float32]().
				Label("Enter price of your job").
				Type("number").
				WrapClass("mt-3").
				Value(&uf.data.PricePerHour).
				PrependIcon("currency_ruble").
				Min(1).
				Step(0.01).
				ID("cost-field"),

			app.If(uf.isRegister, func() app.UI {
				return NewCheckboxField[bool]().
					Label("Remember me").
					Value(&uf.data.IsRemember).
					WrapClass("mt-3").
					ID("remember-field")
			}),

			app.If(uf.isRegister, func() app.UI {
				return NewCheckboxField[bool]().
					Label("I agree to the Terms of Service and Privacy Policy").
					Value(&uf.data.IsAccept).
					WrapClass("mt-1").
					Required(true).
					ID("accept-field")
			}),

			app.If(uf.isRegister, func() app.UI {
				return app.Button().
					Text("Save").
					Class("mt-3 btn btn-primary btn-lg").
					Disabled(!uf.data.Validate(false)).
					OnClick(func(ctx app.Context, e app.Event) {
						if _, err := uf.Ctx.Store.SaveUser(*uf.data.User); err != nil {
							uf.ShowNotification(ctx, "Error", err.Error())

						} else if err := uf.Ctx.Store.LoginByCredentials(ctx, uf.data); err != nil {
							uf.ShowNotification(ctx, "Error", err.Error())
						} else {
							uf.onHide(ctx)
							ctx.Navigate("/index")
						}
					})
			}),

			app.If(uf.isRegister, func() app.UI {
				return app.Raw(`
						<a class="mt-3 icon-link icon-link-hover link-secondary link-underline-light link-underline-opacity-25" href="/login">
						I have already an account: <b>Sign In</b>
							<svg xmlns="http://www.w3.org/2000/svg" class="bi" viewBox="0 0 16 16" aria-hidden="true">
								<path d="M1 8a.5.5 0 0 1 .5-.5h11.793l-3.147-3.146a.5.5 0 0 1 .708-.708l4 4a.5.5 0 0 1 0 .708l-4 4a.5.5 0 0 1-.708-.708L13.293 8.5H1.5A.5.5 0 0 1 1 8z"/>
							</svg>
						</a>
				`)
			}),
		)
}

func (uf *UserForm) Modal() app.UI {
	return NewDialog("userDialog").
		Title("User form").
		Visible(uf.visible).
		IsValid(uf.isValid()).
		OnOk(func(ctx app.Context, e app.Event) {
			if _, err := uf.Ctx.Store.SaveUser(*uf.data.User); err != nil {
				uf.ShowNotification(ctx, "Error save user", err.Error())
			} else {
				uf.onHide(ctx)
			}
		}).
		OnDismiss(func(ctx app.Context, e app.Event) {
			ctx.NewAction("resetUser")
		}).
		Content(uf.El())
}

func (uf *UserForm) onHide(ctx app.Context) {
	ctx.NewAction("hideModal")
	ctx.NewAction("resetUser")
}

func (uf *UserForm) Render() app.UI {
	return app.If(uf.isRegister, uf.El).Else(uf.Modal)
}

func (uf *UserForm) isValid() bool {
	return uf.data.IsValid(false)
}
