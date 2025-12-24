package components

import (
	"accounter/pkg/tools"
	"slices"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

/*
	 type selectValue interface {
		int64 | string
	}
*/
type SelectOption[T comparable] struct {
	Label    string
	Value    T
	Disabled bool
	Hidden   bool
	Selected bool
}

type SelectField[T comparable] struct {
	app.Compo

	id           string
	label        string
	labelClass   string
	wrapperClass string
	prependIcon  string
	Val          *T
	Vals         *[]T
	options      []SelectOption[T]
	clearable    bool
	multiple     bool
	required     bool
}

func NewSelectField[T comparable]() *SelectField[T] {
	return &SelectField[T]{
		Vals: new([]T),
	}
}

func (sf *SelectField[T]) ID(id string) *SelectField[T] {
	sf.id = id
	return sf
}

func (sf *SelectField[T]) Value(v *T) *SelectField[T] {
	sf.Val = v
	return sf
}

func (sf *SelectField[T]) Values(v *[]T) *SelectField[T] {
	sf.Vals = v
	return sf
}

func (sf *SelectField[T]) Clearable(v bool) *SelectField[T] {
	sf.clearable = v
	return sf
}

func (sf *SelectField[T]) Required(v bool) *SelectField[T] {
	sf.required = v
	return sf
}

func (sf *SelectField[T]) Multiple(v bool) *SelectField[T] {
	sf.multiple = v
	return sf
}

func (sf *SelectField[T]) Options(opts []SelectOption[T]) *SelectField[T] {
	sf.options = opts
	return sf
}

func (sf *SelectField[T]) Label(text string) *SelectField[T] {
	sf.label = text

	return sf
}

func (sf *SelectField[T]) LabelClass(cls string) *SelectField[T] {
	sf.labelClass = cls
	return sf
}

func (sf *SelectField[T]) WrappClass(cls string) *SelectField[T] {
	sf.wrapperClass = cls
	return sf
}

func (sf *SelectField[T]) PrependIcon(icon string) *SelectField[T] {
	sf.prependIcon = icon
	return sf
}

func (sf *SelectField[T]) Render() app.UI {
	return NewInputWrapper(sf.id).
		Required(sf.required).
		Label(sf.label).
		LabelClass(sf.labelClass).
		PrependIcon(sf.prependIcon).
		WrapperClass(sf.wrapperClass).
		Clearable(sf.isClearable()).
		Wrap(sf.makeInput())
}

func (sf *SelectField[T]) makeInput() app.UI {
	opts := append(sf.options, SelectOption[T]{
		Value:    *new(T),
		Hidden:   true,
		Disabled: true,
		Selected: true,
	})

	return app.Select().
		ID(sf.id).
		Multiple(sf.multiple).
		Class("form-control").
		OnInput(sf.onInput).
		Body(
			app.Range(opts).Slice(func(i int) app.UI {
				opt := opts[i]

				return app.Option().
					Label(opt.Label).
					Value(opt.Value).
					Disabled(opt.Disabled).
					Hidden(opt.Hidden).
					Selected(sf.isSelected(opt))
			}),
		)
}

func (sf *SelectField[T]) onInput(ctx app.Context, e app.Event) {
	value := ctx.JSSrc().Get("value").String()
	newValue := tools.StringToValue[T](value)

	if sf.multiple {
		for i, v := range *sf.Vals {
			if v == newValue {
				sf.deselect(i)
				return
			}
		}

		*sf.Vals = append(*sf.Vals, newValue)
	} else {
		*sf.Val = newValue
	}
}

func (sf *SelectField[T]) isSelected(o SelectOption[T]) bool {
	if sf.multiple {
		return slices.Contains(*sf.Vals, o.Value)
	}

	res := *sf.Val == o.Value

	return res
}

func (sf *SelectField[T]) isClearable() bool {
	if !sf.clearable {
		return false
	}

	if sf.multiple {
		return len(*sf.Vals) > 0
	}

	return !tools.IsEmpty(sf.Val)
}

func (sf *SelectField[T]) deselect(i int) {
	*sf.Vals = append((*sf.Vals)[:i], (*sf.Vals)[i+1:]...)
}
