package components

import (
	"accounter/pkg/utils"
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Dialog struct {
	app.Compo

	id         string
	title      any
	footer     any
	showHeader bool
	showFooter bool
	size       string
	fullscreen bool
	centered   bool
	backdrop   any

	IsVisible  bool
	IsValidate bool
	Contents   []app.UI

	onOk      func(ctx app.Context, e app.Event)
	onDismiss func(ctx app.Context, e app.Event)
}

func NewDialog(id string) *Dialog {
	return &Dialog{
		id:         id,
		showHeader: true,
		showFooter: true,
		backdrop:   true,
		centered:   true,
		fullscreen: false,
		IsValidate: true,

		onOk:      func(ctx app.Context, e app.Event) {},
		onDismiss: func(ctx app.Context, e app.Event) {},
	}
}

func (d *Dialog) OnMount(ctx app.Context) {
	ctx.Handle("hideModal", func(ctx app.Context, a app.Action) {
		d.Close()
	})
}

func (d *Dialog) Content(c ...app.UI) *Dialog {
	d.Contents = c
	return d
}

func (d *Dialog) IsValid(v bool) *Dialog {
	d.IsValidate = v
	return d
}

func (d *Dialog) OnOk(cb func(ctx app.Context, e app.Event)) *Dialog {
	d.onOk = cb
	return d
}

func (d *Dialog) OnDismiss(cb func(ctx app.Context, e app.Event)) *Dialog {
	d.onDismiss = cb
	return d
}

func (d *Dialog) Persistent(v bool) *Dialog {
	if v {
		d.backdrop = "static"
	} else {
		d.backdrop = true
	}

	return d
}

func (d *Dialog) Fullscreen(v bool) *Dialog {
	d.fullscreen = v
	return d
}

func (d *Dialog) Centered(v bool) *Dialog {
	d.centered = v
	return d
}

func (d *Dialog) Size(s string) *Dialog {
	d.size = s
	return d
}

func (d *Dialog) ShowHeader(v bool) *Dialog {
	d.showHeader = true
	return d
}

func (d *Dialog) ShowFooter(v bool) *Dialog {
	d.showFooter = v
	return d
}

func (d *Dialog) Title(h any) *Dialog {
	d.title = h
	return d
}

func (d *Dialog) Footer(h any) *Dialog {
	d.footer = h
	return d
}

func (d *Dialog) Visible(v bool) *Dialog {
	d.IsVisible = v
	return d
}

func (d *Dialog) getClasses() (result []string) {
	result = append(result, "modal-dialog")

	if !utils.IsEmpty(d.size) {
		result = append(result, fmt.Sprintf("modal-%s", d.size))
	}

	if d.centered {
		result = append(result, "modal-dialog-centered")
	}

	if d.fullscreen {
		result = append(result, "modal-fullscreen")
	}

	return result
}

func (d *Dialog) Close() {
	el := app.Window().
		Get("document").
		Call("querySelector", fmt.Sprintf("#%s .close-modal", d.id))

	if !el.IsNull() && !el.IsUndefined() {
		el.Call("click")
	}
}

func (d *Dialog) Open() {
	el := app.Window().
		Get("document").
		Call("querySelector", fmt.Sprintf("#%s .open-modal", d.id))

	if !el.IsNull() && !el.IsUndefined() {
		el.Call("click")
	}
}

func (d *Dialog) Render() app.UI {
	return app.Div().
		Class("modal").
		Attr("data-bs-backdrop", d.backdrop).
		Attr("data-bs-keyboard", true).
		Attr("aria-labelledby", d.title).
		Attr("aria-hidden", true).
		ID(d.id).
		TabIndex(-1).
		Body(
			app.Div().
				Class(d.getClasses()...).
				Body(
					app.Div().
						Class("modal-content").
						Body(
							app.Button().
								Hidden(true).
								Class("open-modal").
								Attr("data-bs-target", d.id).
								Attr("data-bs-toggle", "modal"),

							// Header
							app.If(d.showHeader, func() app.UI {
								return app.Div().
									Class("modal-header d-flex flex-row justify-content-between").
									Body(
										app.H5().Class("modal-title").Text(d.title),
										NewBtnIcon("close").
											Tooltip("close").
											Attrs("aria-label", "Close").
											Attrs("data-bs-dismiss", "modal").
											BtnClass("close-modal").
											OnClick(d.onDismiss),
									)
							}),

							// Body
							app.If(d.IsVisible, func() app.UI {
								return app.Div().
									Class("modal-body").
									Body(d.Contents...)
							}).Else(func() app.UI {
								return app.Div().Class("modal-body")
							}),

							// Footer
							app.If(d.showFooter, func() app.UI {
								return app.Div().
									Class("modal-footer").
									Body(
										NewBtnIcon("close").
											Text("Close").
											OnClick(func(ctx app.Context, e app.Event) {
												d.Close()
											}),

										NewBtnIcon("save").
											Text("Save").
											Disabled(!d.IsValidate).
											OnClick(func(ctx app.Context, e app.Event) {
												d.onOk(ctx, e)
											}),
									)
							}),
						),
				),
		)
}
