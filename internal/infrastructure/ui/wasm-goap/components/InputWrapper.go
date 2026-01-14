package components

import (
	"accounter/pkg/utils"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type InputWrapper struct {
	app.Compo

	inputID     string
	label       string
	labelClass  string
	wrapClass   string
	prependIcon string
	required    bool
	IsClearable bool

	onClear app.EventHandler
}

func NewInputWrapper(inputID string) *InputWrapper {
	return &InputWrapper{
		inputID: inputID,
		onClear: func(ctx app.Context, e app.Event) {},
	}
}

func (iw *InputWrapper) Label(text string) *InputWrapper {
	iw.label = text
	return iw
}

func (iw *InputWrapper) LabelClass(cls string) *InputWrapper {
	iw.labelClass = cls
	return iw
}

func (iw *InputWrapper) WrapperClass(cls string) *InputWrapper {
	iw.wrapClass = cls
	return iw
}

func (iw *InputWrapper) Required(v bool) *InputWrapper {
	iw.required = v
	return iw
}

func (iw *InputWrapper) Clearable(v bool) *InputWrapper {
	iw.IsClearable = v
	return iw
}

func (iw *InputWrapper) OnClear(cb app.EventHandler) *InputWrapper {
	iw.onClear = cb
	return iw
}

func (iw *InputWrapper) PrependIcon(icon string) *InputWrapper {
	iw.prependIcon = icon
	return iw
}

func (iw *InputWrapper) Wrap(content ...app.UI) app.UI {
	return app.Div().
		Class(iw.wrapClass).
		Body(
			app.If(!utils.IsEmptyValue(iw.label), func() app.UI {
				return NewInputLabel(iw.label, iw.inputID).Required(iw.required).LabelClass(iw.labelClass)
			}),

			app.Div().
				Class("input-group").
				Body(
					app.If(!utils.IsEmptyValue(iw.prependIcon), func() app.UI {
						return app.Div().Class("input-group-prepend").Body(
							NewIcon(iw.prependIcon).Class("input-group-text"),
						)
					}),

					app.Range(content).Slice(func(i int) app.UI {
						return content[i]
					}),

					app.If(iw.IsClearable, func() app.UI {
						return app.Div().Class("input-group-append").Body(
							NewIcon("close").
								Class("input-group-text").
								El().
								Role("button").
								OnClick(iw.onClear),
						)
					}),
				),
		)
}
