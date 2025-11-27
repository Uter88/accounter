package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type Loading struct {
	app.Compo

	spinner string
	Loading bool
}

func NewLoading(loading bool) *Loading {
	return &Loading{
		Loading: loading,
		spinner: "spinner-border",
	}
}

func (l *Loading) Render() app.UI {
	return app.If(l.Loading, func() app.UI {
		return app.Div().
			Style("position", "absolute").
			Style("background", "white").
			Style("opacity", "0.5").
			Style("z-index", "1000").
			Class("d-flex flex-row justify-content-center align-items-center w-100 h-100").
			Body(
				app.Div().
					Class(l.spinner, "text-primary").
					Role("status").
					Body(
						app.Span().Class("visually-hidden").Text("Loading..."),
					),
			)
	}).Else(func() app.UI {
		return app.Span().Hidden(true)
	})
}
