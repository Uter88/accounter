package components

import (
	"accounter/pkg/tools"
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
	err         error
	cols        int
	rows        int

	loaded bool

	formater  func(v string) string
	validator func(ctx app.Context, v string) error
}

func NewInputField[T InputValue]() *InputField[T] {
	return &InputField[T]{
		tp:  "text",
		Val: new(T),
	}
}

func (i *InputField[T]) OnMount(ctx app.Context) {
	ctx.Handle("reset", func(ctx app.Context, a app.Action) {
		i.clear(ctx, app.Event{})
	})
}

func (f *InputField[T]) Render() app.UI {
	if !tools.IsEmpty(*f.Val) {
		f.loaded = true
	}

	if f.tp == "textarea" {
		input := app.Textarea().
			Placeholder(f.placeholder).
			AutoFocus(f.autofocus).
			ReadOnly(f.readonly).
			Class(f.inputClass, f.getValidCls(), "form-control").
			ID(f.id).
			Cols(f.cols).
			Rows(f.rows).
			Text(*f.Val)

		input.OnInput(f.onInput)

		return f.render(input)
	} else {
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

		if err := f.validator(ctx, val.String()); err != nil {
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

func (f *InputField[T]) Validator(fn func(app.Context, string) error) *InputField[T] {
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

func (f *InputField[T]) render(input app.UI) app.UI {
	return NewInputWrapper(f.id).
		Label(f.label).
		LabelClass(f.labelClass).
		WrapperClass(f.wrapClass).
		PrependIcon(f.prependIcon).
		OnClear(f.clear).
		Required(f.required).
		Clearable(f.clearable && !tools.IsEmpty(f.Val)).
		Wrap(input)
}

func (f *InputField[T]) clear(ctx app.Context, e app.Event) {
	ctx.JSSrc().Set("value", "")
	f.onInput(ctx, e)
	app.Window().GetElementByID(f.id).Set("value", *f.Val)
}
