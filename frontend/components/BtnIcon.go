package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type BtnIcon struct {
	app.Compo

	icon      string
	color     string
	iconClass string
	btnClass  string
	btnText   string

	onClick func(ctx app.Context, e app.Event)
}

func NewBtnIcon() *BtnIcon {
	return &BtnIcon{
		icon:      "logout",
		color:     "text-secondary",
		btnText:   "Exit",
		btnClass:  "btn btn-outline-light btn-flat",
		iconClass: "px-1",
	}
}

func (i *BtnIcon) Color(c string) *BtnIcon {
	i.color = c
	return i
}

func (i *BtnIcon) BtnClass(c string) *BtnIcon {
	i.btnClass = c
	return i
}

func (i *BtnIcon) IconClass(c string) *BtnIcon {
	i.iconClass = c
	return i
}

func (i *BtnIcon) Text(c string) *BtnIcon {
	i.btnText = c
	return i
}

func (i *BtnIcon) OnClick(cb func(ctx app.Context, e app.Event)) *BtnIcon {
	i.onClick = cb
	return i
}

func (i *BtnIcon) Render() app.UI {
	icon := NewIcon(i.icon, "").
		Class(i.iconClass)

	btnGroup := app.Div().
		Class("d-flex flex-row align-items-center px-1").
		Body(
			icon,
			app.Button().
				Class(i.btnClass, i.color).
				Text(i.btnText),
		).
		Role("button").
		OnClick(i.onClick)

	return btnGroup
}
