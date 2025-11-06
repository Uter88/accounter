package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Tooltip struct {
	app.Compo

	ttlText     string
	class       string
	ttlPosition string
}

func NewTooltip() *Tooltip {
	return &Tooltip{
		class:       "custom-tooltip-text",
		ttlPosition: "tooltip-left",
	}
}

func (i *Tooltip) TooltipText(c string) *Tooltip {
	i.ttlText = c
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

func (i *Tooltip) TooltipEl() app.HTMLSpan {
	return app.Span().
		Class(i.class, i.ttlPosition).
		Body(
			app.P().Text(i.ttlText),
		)

}

func (i *Tooltip) El() app.HTMLSpan {
	return i.TooltipEl()
}

func (i *Tooltip) Render() app.UI {
	return i.El()
}
