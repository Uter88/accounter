package pages

import (
	"accounter/frontend/common"
	"accounter/frontend/components"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type registrationPage struct {
	app.Compo
	common.BaseComponent
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
				Class("d-flex flex-row justify-content-start mx-auto align-self-center", "login-page").
				Body(
					components.NewUserForm(i.Ctx, true).
						FormClass("login-form"),
				),
		)

}
