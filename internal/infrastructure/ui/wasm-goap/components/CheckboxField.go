package components

import (
	"accounter/pkg/utils"
	"slices"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type CheckboxField[T comparable] struct {
	app.Compo

	id        string
	label     string
	Val       *T
	Vals      *[]T
	required  bool
	readonly  bool
	wrapClass string

	onUpdate app.EventHandler

	loaded bool
}

func NewCheckboxField[T comparable]() *CheckboxField[T] {
	return &CheckboxField[T]{
		Val:      new(T),
		onUpdate: func(ctx app.Context, e app.Event) {},
	}
}

func (f *CheckboxField[T]) Label(text string) *CheckboxField[T] {
	f.label = text
	return f
}

func (f *CheckboxField[T]) ID(id string) *CheckboxField[T] {
	f.id = id
	return f
}

func (f *CheckboxField[T]) OnUpdate(cb app.EventHandler) *CheckboxField[T] {
	f.onUpdate = cb
	return f
}

func (f *CheckboxField[T]) Required(value bool) *CheckboxField[T] {
	f.required = value
	return f
}

func (f *CheckboxField[T]) ReadOnly(value bool) *CheckboxField[T] {
	f.readonly = value
	return f
}

func (f *CheckboxField[T]) WrapClass(cls string) *CheckboxField[T] {
	f.wrapClass = cls
	return f
}

func (f *CheckboxField[T]) Value(v *T) *CheckboxField[T] {
	f.Val = v
	return f
}

func (f *CheckboxField[T]) Values(v *[]T) *CheckboxField[T] {
	f.Vals = v
	return f
}

func (f *CheckboxField[T]) Render() app.UI {
	input := app.Input().
		ReadOnly(f.readonly).
		Checked(f.isChecked()).
		Value(*f.Val).
		Class("form-check-input").
		Type("checkbox")

	input.OnInput(f.onInput)

	return app.Div().Class("form-check", f.wrapClass).Body(

		// Input
		input,

		// Label
		app.If(!utils.IsEmptyValue(f.label), func() app.UI {
			return NewInputLabel(f.label, f.id).LabelClass("form-check-label mx-2")
		}),
	)
}

func (f *CheckboxField[T]) onInput(ctx app.Context, e app.Event) {
	/* 	if !f.loaded {
	   		f.loaded = true
	   	} else if f.required {
	   		if tools.IsEmpty(f.Val) {
	   			input.Class("is-invalid")
	   		} else {
	   			input.Class("is-valid")
	   		}
	   	}
	*/
	value := ctx.JSSrc().Get("value").String()

	switch value {
	case "false":
		value = "true"
	case "true":
		value = "false"
	}

	newValue := utils.StringToValue[T](value)

	if f.Vals != nil {
		needToAdd := true

		for i, v := range *f.Vals {
			if v == newValue {
				f.deselect(i)
				needToAdd = false
				break
			}
		}

		if needToAdd {
			*f.Vals = append(*f.Vals, newValue)
		}
	} else {
		*f.Val = newValue
	}

	f.onUpdate(ctx, e)
}

func (f *CheckboxField[T]) isChecked() bool {
	if f.Vals != nil {
		return slices.Contains(*f.Vals, *f.Val)
	}

	return !utils.IsEmpty(f.Val)
}

func (sf *CheckboxField[T]) deselect(i int) {
	*sf.Vals = append((*sf.Vals)[:i], (*sf.Vals)[i+1:]...)
}
