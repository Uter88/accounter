package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type BtnIcon struct {
	app.Compo

	icon       string
	color      string
	iconClass  string
	btnClasses []string
	btnText    string
	tooltip    string
	attrs      map[string]any
	role       string
	IsDisabled bool

	onClick func(ctx app.Context, e app.Event)
}

func NewBtnIcon(icon string) *BtnIcon {
	return &BtnIcon{
		icon:       icon,
		color:      "text-secondary",
		btnClasses: []string{"btn btn-outline-light btn-flat"},
		iconClass:  "px-1",
		attrs:      make(map[string]any),
		role:       "button",
		onClick:    func(ctx app.Context, e app.Event) {},
	}
}

func (i *BtnIcon) Target(id string) *BtnIcon {
	i.attrs["data-bs-target"] = id
	i.attrs["data-bs-toggle"] = "modal"
	return i
}

func (i *BtnIcon) Attrs(k string, v any) *BtnIcon {
	i.attrs[k] = v
	return i
}

func (i *BtnIcon) Tooltip(c string) *BtnIcon {
	i.tooltip = c
	return i
}

func (i *BtnIcon) Color(c string) *BtnIcon {
	i.color = c
	return i
}

func (i *BtnIcon) Disabled(v bool) *BtnIcon {
	i.IsDisabled = v
	return i
}

func (i *BtnIcon) BtnClass(c string) *BtnIcon {
	i.btnClasses = append(i.btnClasses, c)
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
	icon := NewIcon(i.icon).
		Tooltip(i.tooltip).
		Class(i.iconClass)

	btnGroup := app.Button().
		Type("button").
		Attr("data-bs-toggle", "button").
		Role(i.role).
		Class("d-flex flex-row align-items-center px-1").
		Class(i.btnClasses...).
		Class(i.color).
		Disabled(i.IsDisabled).
		Body(
			icon,
			app.Text(i.btnText),
		).OnClick(func(ctx app.Context, e app.Event) {
		if !i.IsDisabled {
			i.onClick(ctx, e)
		}
	})

	for k, v := range i.attrs {
		btnGroup.Attr(k, v)
	}

	return btnGroup
}
