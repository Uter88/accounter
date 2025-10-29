package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Icon struct {
	app.Compo

	name  string
	color string
	class string
}

func NewIcon(name string) *Icon {
	return &Icon{name: name, color: "text-secondary"}
}

func (i *Icon) Color(c string) *Icon {
	i.color = c
	return i
}

func (i *Icon) Class(c string) *Icon {
	i.class = c
	return i
}

func (i *Icon) El() app.HTMLSpan {
	return app.Span().
		Class("material-symbols-outlined", i.color, i.class).
		Text(i.name)
}

func (i *Icon) Render() app.UI {
	return i.El()
}
