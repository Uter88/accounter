package common

import (
	"accounter/config"
	"accounter/internal/infrastructure/ui/wasm-goap/models"
	"accounter/internal/infrastructure/ui/wasm-goap/store"
	"accounter/pkg/logger"

	"context"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type BaseComponent struct {
	Ctx AppContext

	notificationPermission app.NotificationPermission
}

func NewBaseComponent(ctx AppContext) BaseComponent {
	return BaseComponent{
		Ctx: ctx,
	}
}

func (h *BaseComponent) OnMount(ctx app.Context) {
	h.notificationPermission = ctx.Notifications().Permission()
}

func (h *BaseComponent) EnableNotifications(ctx app.Context) {
	h.notificationPermission = ctx.Notifications().RequestPermission()
}

func (h *BaseComponent) ShowNotification(ctx app.Context, title, msg string) {
	ctx.Notifications().New(app.Notification{
		Title: title,
		Body:  msg,
		Icon:  "/web/icons/logo.svg",
	})
}

func (h *BaseComponent) ShowConfirm(ctx app.Context, title string) bool {
	confirm := app.Window().Get("confirm")
	result := confirm.Invoke(title)

	return result.Bool()
}

type AppContext struct {
	context.Context
	Store     *store.Store
	Logger    logger.Logger
	Config    config.Config
	Websocket *models.WebsocketClient
}
