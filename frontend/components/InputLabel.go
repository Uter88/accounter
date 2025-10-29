package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type InputLabel struct {
	app.Compo

	inputID  string
	text     string
	required bool
}

func NewInputLabel(text, inputID string) *InputLabel {
	return &InputLabel{text: text}
}

func (l *InputLabel) Required(v bool) *InputLabel {
	l.required = v
	return l
}

func (l *InputLabel) Render() app.UI {
	text := l.text

	if l.required {
		text = "*" + text
	}

	return app.Label().
		For(l.inputID).
		Text(text).
		Class("text-secondary")
}
