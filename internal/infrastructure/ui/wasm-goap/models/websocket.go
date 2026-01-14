package models

import (
	"accounter/config"
	"accounter/internal/domain/shared"
	"accounter/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WebsocketClient struct {
	url                  string
	conn                 app.Value
	logger               logger.Logger
	reconnectionInterval time.Duration
	pingInterval         time.Duration
}

func NewWebsocket(config config.Config, logger logger.Logger) WebsocketClient {
	return WebsocketClient{
		url:                  config.HTTP.Websocket.URL,
		logger:               logger,
		reconnectionInterval: config.HTTP.Websocket.ReconnectionInterval,
		pingInterval:         config.HTTP.Websocket.PingInterval,
	}
}

func (w *WebsocketClient) Connect(token string) {
	w.logger.Info("Connecting to websocket")

	ws := app.Window().Get("WebSocket")
	w.conn = ws.New(w.url)

	w.conn.Set("onopen", app.FuncOf(func(this app.Value, args []app.Value) any {
		w.logger.Info("Websocket connected")
		return w.SendMessage(shared.WsMessageAuth, token)
	}))

	w.conn.Set("onclose", app.FuncOf(func(this app.Value, args []app.Value) any {
		w.logger.Warn("Webscoket connection closed. Reconnecting...")

		time.AfterFunc(w.reconnectionInterval, func() {
			w.Connect(token)
		})

		return nil
	}))

	w.conn.Set("onmessage", app.FuncOf(func(this app.Value, args []app.Value) any {
		if err := w.onMessage(args); err != nil {
			w.logger.Errorf("Websocket: %s", err.Error())
		}

		return this
	}))
}

func (w *WebsocketClient) onMessage(args []app.Value) error {
	msg, err := w.parseMessage(args)

	if err != nil {
		return err
	}

	w.logger.Infof("Websocket: new message: %s", msg)

	switch msg.Type {
	case shared.WsMessageAuth:
		if msg.Error != "" {
			return fmt.Errorf("authorization error: %s", msg.Error)
		}
	case shared.WsMessagePong:
		w.DoPing()
	case shared.WsMessagePing:
		w.DoPong()
	}

	return nil
}

func (w *WebsocketClient) parseMessage(args []app.Value) (msg shared.WebsocketMessage, err error) {
	if len(args) == 0 {
		err = errors.New("empty message")
		return
	}

	data := args[0].Get("data")

	if data.Type() == app.TypeString {
		err = json.Unmarshal([]byte(data.String()), &msg)
	} else {
		err = fmt.Errorf("invalid websocket data type: %s", data.Type())
	}

	return
}

func (w *WebsocketClient) DoPing() {
	time.AfterFunc(w.pingInterval, func() {
		w.SendMessage(shared.WsMessagePing, nil)
	})
}

func (w *WebsocketClient) DoPong() {
	w.SendMessage(shared.WsMessagePong, nil)
}

func (w *WebsocketClient) SendMessage(msgType shared.WsMessageType, payload any) any {
	if w.conn == nil {
		w.logger.Warn("Websocket connection is nil")
		return nil
	}

	msg := shared.WebsocketMessage{
		Type:    msgType,
		Payload: payload,
	}

	return w.conn.Call("send", msg.ToString())
}
