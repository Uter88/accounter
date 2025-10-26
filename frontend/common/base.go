package common

import (
	"accounter/config"
	"accounter/frontend/store"
	"context"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type BaseComponent struct {
	app.Compo
	Ctx AppContext

	notificationPermission app.NotificationPermission
}

func NewBaseComponent(ctx AppContext) BaseComponent {
	return BaseComponent{Ctx: ctx}
}

func (h *BaseComponent) OnMount(ctx app.Context) {
	h.notificationPermission = ctx.Notifications().Permission()
}

func (h *BaseComponent) EnableNotifications(ctx app.Context, e app.Event) {
	h.notificationPermission = ctx.Notifications().RequestPermission()
}

func (h *BaseComponent) ShowNotification(ctx app.Context, title, msg string) {
	ctx.Notifications().New(app.Notification{
		Title: title,
		Body:  msg,
	})
}

type AppContext struct {
	context.Context
	Store  *store.Store
	Logger config.Logger
	Config config.Config
}
