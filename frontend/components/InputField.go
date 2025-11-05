package components

import (
	"accounter/tools"
	"errors"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type InputValue interface {
	int | float64 | float32 | string
}

type InputField[T InputValue] struct {
	app.Compo

	id          string
	inputClass  string
	wrapClass   string
	labelClass  string
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
	readonly    bool
	pattern     string
	validCls    string
	err         error

	loaded bool

	formater  func(v string) string
	validator func(v string) error
}

func NewInputField[T InputValue]() *InputField[T] {
	return &InputField[T]{tp: "text"}
}

func (f *InputField[T]) Render() app.UI {
	input := app.Input().
		Type(f.tp).
		ID(f.id).
		Class(f.inputClass, f.getValidCls(), "form-control").
		Value(*f.Val).
		Placeholder(f.placeholder).
		AutoFocus(f.autofocus).
		ReadOnly(f.readonly)

	if f.tp == "number" {
		if !tools.IsEmpty(f.step) {
			input.Step(f.step)
		}

		if !tools.IsEmpty(f.min) {
			input.Min(f.min)
		}

		if !tools.IsEmpty(f.max) {
			input.Max(f.max)
		}
	}

	if !tools.IsEmpty(f.pattern) {
		input.Pattern(f.pattern)
	}

	input.OnInput(f.onInput)

	return f.render(input)
}

func (f *InputField[T]) onInput(ctx app.Context, e app.Event) {
	if f.formater != nil {
		val := ctx.JSSrc().Get("value")
		ctx.JSSrc().Set("value", f.formater(val.String()))
	}

	h := f.ValueTo(f.Val)
	h(ctx, e)

	f.err = nil

	if !f.loaded {
		f.loaded = true
	} else if f.required {
		if tools.IsEmpty(f.Val) {
			f.err = errors.New("required field")
		}
	}

	if f.validator != nil {
		val := ctx.JSSrc().Get("value")

		if err := f.validator(val.String()); err != nil {
			f.err = err
		}
	}
}

func (f *InputField[T]) getValidCls() string {
	if f.err != nil {
		return "is-invalid"
	}

	if f.loaded && (f.required || f.validator != nil) {
		return "is-valid"
	}

	return ""
}

func (f *InputField[T]) Formater(fn func(string) string) *InputField[T] {
	f.formater = fn
	return f
}

func (f *InputField[T]) Validator(fn func(string) error) *InputField[T] {
	f.validator = fn
	return f
}

func (f *InputField[T]) Pattern(p string) *InputField[T] {
	f.pattern = p
	return f
}

func (f *InputField[T]) ID(id string) *InputField[T] {
	f.id = id
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

func (f *InputField[T]) InputClass(c string) *InputField[T] {
	f.inputClass = c
	return f
}

func (f *InputField[T]) LabelClass(c string) *InputField[T] {
	f.labelClass = c
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

func (f *InputField[T]) render(input app.HTMLInput) app.UI {
	return app.Div().Class(f.wrapClass).Body(

		// Label
		app.If(!tools.IsEmptyValue(f.label), func() app.UI {
			return NewInputLabel(f.label, f.id).Required(f.required).LabelClass(f.labelClass)
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
						NewIcon("close").
							Class("input-group-text").
							El().
							Role("button").
							OnClick(func(ctx app.Context, e app.Event) {
								ctx.JSSrc().Set("value", "")
								f.onInput(ctx, e)
								app.Window().GetElementByID(f.id).Set("value", *f.Val)
							}),
					)
				}),

				// app.Div().Class("invalid-feedback").Text("Invalid"),
			),
	)
}
