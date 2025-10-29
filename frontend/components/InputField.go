package components

import (
	"accounter/tools"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type InputValue interface {
	int | float64 | string
}

type InputField[T InputValue] struct {
	app.Compo

	id          string
	inputClass  string
	wrapClass   string
	tp          string
	label       string
	placeholder string
	value       T
	rcv         *T
	autofocus   bool
	required    bool
	clearable   bool
	min         any
	max         any
	step        float64
	prependIcon string

	onInput app.EventHandler
}

func NewInputField[T InputValue]() *InputField[T] {
	return &InputField[T]{tp: "text"}
}

func (f *InputField[T]) Render() app.UI {
	input := app.Input().
		Type(f.tp).
		ID(f.id).
		Class("form-control", f.inputClass).
		Value(f.value).
		Min(f.min).
		Max(f.max).
		Step(f.step).
		Placeholder(f.placeholder).
		AutoFocus(f.autofocus).
		OnInput(f.onInput)

	return app.Div().Class(f.wrapClass).Body(

		// Label
		app.If(!tools.IsEmptyValue(f.label), func() app.UI {
			return NewInputLabel(f.label, f.id).Required(f.required)
		}),

		// Content
		app.Div().
			Class("input-group").
			Body(
				app.If(!tools.IsEmptyValue(f.prependIcon), func() app.UI {
					return app.Div().Class("input-group-prepend").Body(
						NewIcon(f.prependIcon).Class("input-group-text"),
					)
				}),

				// Input
				input,

				// Clear icon
				app.If(f.clearable && !tools.IsEmpty(f.value), func() app.UI {
					return app.Div().Class("input-group-append").Body(
						NewIcon("close").Class("input-group-text").El().OnClick(func(ctx app.Context, e app.Event) {
							input.JSValue().Set("value", *new(T))
						}),
					)
				}),
			),
	)
}

func (f *InputField[T]) ID(id string) *InputField[T] {
	f.id = id
	return f
}

func (f *InputField[T]) InputClass(c string) *InputField[T] {
	f.inputClass = c
	return f
}

func (f *InputField[T]) WrapClass(c string) *InputField[T] {
	f.wrapClass = c
	return f
}

func (f *InputField[T]) Autofocus(v bool) *InputField[T] {
	f.autofocus = v
	return f
}

func (f *InputField[T]) Clearable(v bool) *InputField[T] {
	f.clearable = v
	return f
}

func (f *InputField[T]) Required(v bool) *InputField[T] {
	f.required = v
	return f
}

func (f *InputField[T]) Type(tp string) *InputField[T] {
	f.tp = tp
	return f
}

func (f *InputField[T]) Label(text string) *InputField[T] {
	f.label = text
	return f
}

func (f *InputField[T]) Placeholder(text string) *InputField[T] {
	f.placeholder = text
	return f
}

func (f *InputField[T]) Value(value T) *InputField[T] {
	f.value = value
	return f
}

func (f *InputField[T]) Rcv(value *T) *InputField[T] {
	f.rcv = value
	return f
}

func (f *InputField[T]) Min(value any) *InputField[T] {
	f.min = value
	return f
}

func (f *InputField[T]) Max(value any) *InputField[T] {
	f.max = value
	return f
}

func (f *InputField[T]) Step(value float64) *InputField[T] {
	f.step = value
	return f
}

func (f *InputField[T]) PrependIcon(value string) *InputField[T] {
	f.prependIcon = value
	return f
}

func (f *InputField[T]) OnInput(h app.EventHandler) *InputField[T] {
	f.onInput = h
	return f
}
