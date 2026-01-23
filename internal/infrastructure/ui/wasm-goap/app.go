package wasmgoap

import (
	"accounter/config"
	"accounter/internal/infrastructure/ui/wasm-goap/common"
	"accounter/internal/infrastructure/ui/wasm-goap/models"
	"accounter/internal/infrastructure/ui/wasm-goap/pages"
	"accounter/internal/infrastructure/ui/wasm-goap/store"
	"accounter/pkg/logger"

	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type App struct {
	ctx common.AppContext
}

func NewApp(config config.Config, logger logger.Logger) App {
	wsClient := models.NewWebsocket(config, logger)
	store := store.NewStore(config, &wsClient)

	return App{
		ctx: common.AppContext{
			Config:    config,
			Logger:    logger,
			Store:     store,
			Websocket: &wsClient,
		},
	}
}

func (a *App) Name() string {
	return fmt.Sprintf("HTTP client on 0.0.0.0:%d", a.ctx.Config.Client.Port)
}

func (a *App) Run(ctx context.Context) error {
	app.Route("/", func() app.Composer { return pages.NewDefaultPage(a.ctx) })
	app.Route("/login", func() app.Composer { return pages.NewLoginPage(a.ctx) })
	app.Route("/registration", func() app.Composer { return pages.NewRegistrationPage(a.ctx) })
	app.Route("/index", func() app.Composer { return pages.NewIndexPage(a.ctx) })

	app.RunWhenOnBrowser()

	handler := app.Handler{
		Name: "Accounter",
		Icon: app.Icon{
			Default: "/web/icons/favorite.png",
			SVG:     "/web/icons/logo-bg.svg",
		},
		Title:       "AccApp",
		Description: "Accounter application",
		Styles: []string{
			"/web/styles/bootstrap.min.css",
			//"https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200",
			"/web/styles/base.css",
			"/web/styles/popups.css",
			"/web/styles/components.css",
		},
		//Fonts: []string{"/web/fonts/material.woff2"},

		Scripts: []string{
			"/web/js/main.js",
			"/web/js/echarts/echarts.min.js",
			"/web/js/bootstrap.bundle.min.js",
		},
		//Resources: app.LocalDir(".."),
	}

	serv := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", a.ctx.Config.Client.Port),
		ReadTimeout:  a.ctx.Config.HTTP.ReadTimeout,
		WriteTimeout: a.ctx.Config.HTTP.WriteTimeout,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
		Handler: &handler,
	}

	return serv.ListenAndServe()
}
