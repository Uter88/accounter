package frontend

import (
	"accounter/config"
	"accounter/frontend/common"
	"accounter/frontend/pages"
	"accounter/frontend/store"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type App struct {
	ctx common.AppContext
}

func NewApp(config config.Config, logger config.Logger) App {
	return App{
		ctx: common.AppContext{
			Config: config,
			Logger: logger,
			Store:  store.NewStore(config),
		},
	}
}

func (a *App) Run(ctx context.Context) error {
	app.Route("/login", func() app.Composer { return pages.NewLoginPage(a.ctx) })
	app.Route("/index", func() app.Composer { return pages.NewIndexPage(a.ctx) })
	app.RunWhenOnBrowser()

	serv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", a.ctx.Config.Client.Port),
		ReadTimeout:  time.Minute * 20,
		WriteTimeout: time.Minute * 20,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
		Handler: &app.Handler{
			Name: "Accounter",
			Icon: app.Icon{
				Default: "/web/icons/favorite.png",
				SVG:     "/web/icons/favorite.svg",
			},
			Title:       "Accounter application",
			Description: "Accounter application",
			Styles: []string{
				"https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css",
				"https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200",
				"/web/styles.css",
			},
			Scripts: []string{
				"https://cdn.jsdelivr.net/npm/bootstrap@5.3.8/dist/js/bootstrap.bundle.min.js",
			},
		},
	}

	a.ctx.Logger.Infof("Start client HTTP server on %d port", a.ctx.Config.Client.Port)

	return serv.ListenAndServe()
}
