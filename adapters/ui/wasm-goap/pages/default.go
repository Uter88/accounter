package pages

import (
	"accounter/adapters/ui/wasm-goap/common"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type defaultPage struct {
	app.Compo
	common.BaseComponent
}

func NewDefaultPage(ctx common.AppContext) *defaultPage {
	return &defaultPage{
		BaseComponent: common.NewBaseComponent(ctx),
	}
}

func (dp *defaultPage) OnMount(ctx app.Context) {
	ctx.Navigate("/index")
}

func (dp *defaultPage) Render() app.UI {
	return app.Div()
}
