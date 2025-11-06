package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Icon struct {
	app.Compo

	name         string
	color        string
	class        string
	tooltip      string
	alignTooltip string
}

func NewIcon(name, tooltip string) *Icon {
	return &Icon{
		name:         name,
		color:        "text-secondary",
		tooltip:      tooltip,
		alignTooltip: "tooltip-bottom",
	}
}

func (i *Icon) Color(c string) *Icon {
	i.color = c
	return i
}

func (i *Icon) Class(c string) *Icon {
	i.class = c
	return i
}

func (i *Icon) AlignClass(c string) *Icon {
	i.alignTooltip = c
	return i
}

func (i *Icon) El() app.HTMLSpan {
	if i.tooltip != "" {
		return app.Span().
			Class("custom-tooltip").
			Body(
				app.Span().
					Class("material-symbols-outlined", i.color, i.class).
					Text(i.name),

				NewTooltip().
					TooltipText(i.tooltip).
					AlignClass(i.alignTooltip),
			)

	} else {
		return app.Span().
			Class("material-symbols-outlined", i.color, i.class).
			Text(i.name)
	}

}

func (i *Icon) Render() app.UI {
	return i.El()
}
