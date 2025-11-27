package components

import (
	"accounter/pkg/tools"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type inputType string

const (
	InputTypeDate          inputType = "date"
	InputTypeTime          inputType = "time"
	InputTypeDatetime      inputType = "datetime"
	InputTypeDateRange     inputType = "daterange"
	InputTypeTimeRange     inputType = "timerange"
	InputTypeDatetimeRange inputType = "datetimerange"
)

type Calendar struct {
	app.Compo

	id          string
	tp          inputType
	wrapClass   string
	labelClass  string
	label       string
	prependIcon string
	required    bool
	clearable   bool
	min         int64
	max         int64
	Val         *int64
	Vals        []*int64

	onUpdate app.EventHandler
	updates  map[int]bool
}

func (f *Calendar) PrependIcon(value string) *Calendar {
	f.prependIcon = value
	return f
}

func NewCalendar() *Calendar {
	return &Calendar{
		tp:       InputTypeDate,
		onUpdate: func(ctx app.Context, e app.Event) {},
		updates:  make(map[int]bool),
	}
}

func (c *Calendar) Required(v bool) *Calendar {
	c.required = v
	return c
}

func (c *Calendar) Min(m int64) *Calendar {
	c.min = m
	return c
}

func (c *Calendar) Max(m int64) *Calendar {
	c.max = m
	return c
}

func (c *Calendar) Clearable(v bool) *Calendar {
	c.clearable = v
	return c
}

func (c *Calendar) ID(id string) *Calendar {
	c.id = id
	return c
}

func (c *Calendar) WrapClass(cls string) *Calendar {
	c.wrapClass = cls
	return c
}

func (c *Calendar) LabelClass(cls string) *Calendar {
	c.labelClass = cls
	return c
}

func (c *Calendar) Value(v *int64) *Calendar {
	c.Val = v
	return c
}

func (c *Calendar) Values(start, stop *int64) *Calendar {
	c.Vals = []*int64{start, stop}
	c.updates[0] = false
	c.updates[1] = false

	return c
}

func (c *Calendar) OnUpdate(cb app.EventHandler) *Calendar {
	c.onUpdate = cb
	return c
}

func (c *Calendar) Type(tp inputType) *Calendar {
	c.tp = tp
	return c
}

func (c *Calendar) Label(text string) *Calendar {
	c.label = text
	return c
}

func (c *Calendar) Render() app.UI {
	return NewInputWrapper(c.id).
		WrapperClass(c.wrapClass).
		LabelClass(c.labelClass).
		Required(c.required).
		Label(c.label).
		PrependIcon(c.prependIcon).
		Clearable(c.clearable && !tools.IsEmpty(c.Val)).
		Wrap(c.makeInputs()...)
}

func (c *Calendar) makeInputs() (inputs []app.UI) {
	values := make([]*int64, 0)

	if len(c.Vals) > 0 {
		values = append(values, c.Vals[0], c.Vals[1])
	} else {
		values = append(values, c.Val)
	}

	for i, v := range values {
		minDate, maxDate := c.getMinMax(i, values)

		input := app.Input().
			Type(c.getType()).
			Min(minDate).
			Max(maxDate).
			Value(c.formatValue(*v)).
			Class("form-control").
			OnInput(func(ctx app.Context, e app.Event) {
				c.onInput(ctx, e, i, v)
			})

		inputs = append(inputs, input)
	}

	return
}

func (c *Calendar) onInput(ctx app.Context, e app.Event, index int, dest *int64) {
	value := ctx.JSSrc().Get("value")
	dt, err := time.ParseInLocation(c.getFormat(), value.String(), time.Local)

	if err == nil {
		if index == 1 && c.tp == InputTypeDateRange {
			dt = time.Date(dt.Year(), dt.Month(), dt.Day(), 23, 59, 59, 0, dt.Location())
		}

		ts := dt.Unix()
		*dest = ts

		c.updates[index] = true

		if c.isReadyToUpdate() {
			c.onUpdate(ctx, e)
			c.resetUpdates()
		}
	}
}

func (c *Calendar) resetUpdates() {
	for i := range c.updates {
		c.updates[i] = false
	}
}

func (c *Calendar) isReadyToUpdate() bool {
	for _, u := range c.updates {
		if !u {
			return false
		}
	}

	return true
}

func (c *Calendar) formatValue(v int64) string {
	format := c.getFormat()
	return time.Unix(v, 0).Format(format)
}

func (c *Calendar) getFormat() string {
	switch c.tp {
	case InputTypeDate, InputTypeDateRange:
		return time.DateOnly

	case InputTypeTime, InputTypeTimeRange:
		return "15:04"

	case InputTypeDatetime, InputTypeDatetimeRange:
		return "2006-01-02T15:04"
	}

	return ""
}

func (c *Calendar) getType() string {
	switch c.tp {
	case InputTypeDate, InputTypeDateRange:
		return "date"

	case InputTypeTime, InputTypeTimeRange:
		return "time"

	case InputTypeDatetime, InputTypeDatetimeRange:
		return "datetime-local"
	}

	return ""
}

func (c *Calendar) getMinMax(index int, values []*int64) (minDate, maxDate string) {
	format := c.getFormat()
	var minTs, maxTs int64

	if len(values) > 1 {
		if index == 0 {
			maxTs = *values[1]
		} else {
			minTs = *values[0]
		}
	}

	if c.min > 0 {
		minTs = max(c.min, minTs)
	}

	if c.max > 0 {
		maxTs = min(c.max, maxTs)
	}

	if minTs > 0 {
		minDate = time.Unix(minTs, 0).Format(format)
	}

	if maxTs > 0 {
		maxDate = time.Unix(maxTs, 0).Format(format)
	}

	return
}
