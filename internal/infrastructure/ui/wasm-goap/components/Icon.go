package components

import (
	"accounter/pkg/utils"

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

func NewIcon(name string) *Icon {
	return &Icon{
		name:         name,
		color:        "text-secondary",
		alignTooltip: "tooltip-bottom",
	}
}

func (i *Icon) Tooltip(t string) *Icon {
	i.tooltip = t
	return i
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
	icon := app.Span().
		Class("material-symbols-outlined", i.color, i.class).
		Text(i.name)

	return app.Span().Body(
		app.If(!utils.IsEmpty(i.tooltip), func() app.UI {
			return NewTooltip().Text(i.tooltip).AlignClass(i.alignTooltip).El(icon)
		}).Else(func() app.UI {
			return icon
		}),
	)
}

func (i *Icon) Render() app.UI {
	return i.El()
}
