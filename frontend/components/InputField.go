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
	Val         *T
	autofocus   bool
	required    bool
	clearable   bool
	min         any
	max         any
	step        float64
	prependIcon string
}

func NewInputField[T InputValue]() *InputField[T] {
	return &InputField[T]{tp: "text"}
}

func (f *InputField[T]) Render() app.UI {
	input := app.Input().
		Type(f.tp).
		ID(f.id).
		Class("form-control", f.inputClass).
		Value(*f.Val).
		Min(f.min).
		Max(f.max).
		Step(f.step).
		Placeholder(f.placeholder).
		AutoFocus(f.autofocus).
		OnInput(f.ValueTo(f.Val))

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
				app.If(f.clearable && !tools.IsEmpty(f.Val), func() app.UI {
					return app.Div().Class("input-group-append").Body(
						NewIcon("close").Class("input-group-text").El().Role("button").OnClick(func(ctx app.Context, e app.Event) {
							*f.Val = *new(T)
							app.Window().GetElementByID(f.id).Set("value", *f.Val)
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

func (f *InputField[T]) Value(value *T) *InputField[T] {
	f.Val = value
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
