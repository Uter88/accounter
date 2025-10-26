package layouts

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type MainLayout struct {
	app.Compo
	child app.Composer
}

func (l *MainLayout) Render() app.UI {
	return app.Div().Body(
		app.Div().Text("Header"),
		app.Div().Body(l.child),
		app.Div().Text("Footer"),
	)
}
