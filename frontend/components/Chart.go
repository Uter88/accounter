package components

import (
	"accounter/frontend/models"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Chart struct {
	app.Compo

	id      string
	options models.Object
	width   string
	height  string
	el      app.Value
	lazy    bool

	onClick func(ctx app.Context, params models.Params)
}

func (c *Chart) OnMount(ctx app.Context) {
	if !c.lazy {
		c.Draw(ctx)
	}
}

func (c *Chart) OnResize(ctx app.Context) {
	c.Resize()
}

func (c *Chart) Resize() *Chart {
	if c.isValid() {
		c.el.Call("resize")
	}

	return c
}

func (c *Chart) isValid() bool {
	if c.el == nil {
		return false
	}

	if c.el.IsUndefined() || c.el.IsNull() {
		return false
	}

	return true
}

func (c *Chart) Draw(ctx app.Context) *Chart {
	if c.isValid() {
		c.el.Call("clear")
	}

	c.el = app.Window().
		Get("echarts").
		Call("init", app.Window().GetElementByID(c.id))

	c.el.Call("setOption", c.options)

	if c.onClick != nil {
		c.el.Call("on", "click", "series", app.FuncOf(func(this app.Value, args []app.Value) any {
			params := models.Params{
				SeriesType: args[0].Get("seriesType").String(),
			}
			c.onClick(ctx, params)

			return nil
		}))
	}

	return c
}

func (c *Chart) OnClick(fn func(ctx app.Context, params models.Params)) *Chart {
	c.onClick = fn
	return c
}

func (c *Chart) Lazy(v bool) *Chart {
	c.lazy = v
	return c
}

func (c *Chart) Width(w string) *Chart {
	c.width = w
	return c
}

func (c *Chart) Height(h string) *Chart {
	c.height = h
	return c
}

func (c *Chart) SetOptions(options *models.ChartOptions) *Chart {
	if options != nil {
		c.options = options.Encode()
	}

	return c
}

func (c *Chart) SetRawOptions(options models.Object) *Chart {
	c.options = options
	return c
}

func (c *Chart) Render() app.UI {
	return app.Div().
		ID(c.id).
		Style("width", c.width).
		Style("height", c.height).
		Style("display", "flex").
		Style("justify-content", "center").
		Style("align-items", "center")
}

func NewChart(id string) *Chart {
	return &Chart{
		id:      id,
		options: make(models.Object),
		width:   "100%",
		height:  "100%",
	}
}
