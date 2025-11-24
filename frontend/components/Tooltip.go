package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Tooltip struct {
	app.Compo

	text        string
	class       string
	ttlPosition string
}

func NewTooltip() *Tooltip {
	return &Tooltip{
		class:       "custom-tooltip-text",
		ttlPosition: "tooltip-left",
	}
}

func (i *Tooltip) Text(c string) *Tooltip {
	i.text = c
	return i
}

func (i *Tooltip) Class(c string) *Tooltip {
	i.class = c
	return i
}

func (i *Tooltip) AlignClass(c string) *Tooltip {
	i.ttlPosition = c
	return i
}

func (i *Tooltip) El(contents ...app.UI) app.HTMLSpan {
	return app.Span().
		Class("custom-tooltip").
		Body(
			app.Range(contents).Slice(func(i int) app.UI {
				return contents[i]
			}),

			app.Span().
				Class(i.class, i.ttlPosition).
				Body(
					app.P().Text(i.text),
				),
		)

}

func (i *Tooltip) Render() app.UI {
	return i.El()
}
