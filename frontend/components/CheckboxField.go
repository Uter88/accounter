package components

import (
	"accounter/tools"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type CheckboxField struct {
	app.Compo

	id        string
	label     string
	Val       *bool
	required  bool
	readonly  bool
	wrapClass string

	loaded bool
}

func NewCheckboxField() *CheckboxField {
	return &CheckboxField{}
}

func (f *CheckboxField) Label(text string) *CheckboxField {
	f.label = text
	return f
}

func (f *CheckboxField) ID(id string) *CheckboxField {
	f.id = id
	return f
}

func (f *CheckboxField) Checked(value *bool) *CheckboxField {
	f.Val = value
	return f
}

func (f *CheckboxField) Required(value bool) *CheckboxField {
	f.required = value
	return f
}

func (f *CheckboxField) ReadOnly(value bool) *CheckboxField {
	f.readonly = value
	return f
}

func (f *CheckboxField) WrapClass(cls string) *CheckboxField {
	f.wrapClass = cls

	return f
}

func (f *CheckboxField) Render() app.UI {
	input := app.Input().
		ReadOnly(f.readonly).
		Checked(*f.Val).
		Value("").
		Class("form-check-input").
		Type("checkbox")

	input.OnInput(f.onInput(input))

	return app.Div().Class("form-check", f.wrapClass).Body(

		// Input
		input,

		// Label
		app.If(!tools.IsEmptyValue(f.label), func() app.UI {
			return NewInputLabel(f.label, f.id).LabelClass("form-check-label mx-2")
		}),
	)
}

func (f *CheckboxField) onInput(input app.HTMLInput) app.EventHandler {
	if !f.loaded {
		f.loaded = true
	} else if f.required {
		if tools.IsEmpty(f.Val) {
			input.Class("is-invalid")
		} else {
			input.Class("is-valid")
		}
	}

	return f.ValueTo(f.Val)
}
