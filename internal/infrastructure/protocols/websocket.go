package protocols

import (
	"accounter/config"
	"accounter/internal/domain/event"
	"accounter/internal/domain/shared"
	"accounter/internal/domain/user"
	"accounter/pkg/logger"
	"accounter/pkg/tools"
	"context"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type (
	// Websocket service
	websocketService struct {
		websocket.Upgrader

		// Authorization service
		authService authService

		// Application config
		config config.Config

		logger logger.Logger

		// Connections store
		connections tools.SyncMap[int64, []*websocketConn]
	}

	// Websocket connection
	websocketConn struct {

		// Request context
		ctx context.Context

		// Mutex and connection
		mu   sync.Mutex
		conn *websocket.Conn

		// Read and write deadlines
		readDeadline  time.Duration
		writeDeadline time.Duration

		// Authorizated user instance
		user user.CurrentUser
	}

	// Authorization service
	authService interface {
		LoginByToken(ctx context.Context, token string, cfg config.Config) (result user.CurrentUser, err error)
		LoginByCredentials(ctx context.Context, login, password string, cfg config.Config) (result user.CurrentUser, err error)
	}
)

// NewWebsocketService creates new websocketService instance
func NewWebsocketService(authService authService, cfg config.Config, logger logger.Logger) *websocketService {
	return &websocketService{
		authService: authService,
		config:      cfg,
		logger:      logger,

		connections: tools.NewSyncMap[int64, []*websocketConn](),

		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Name of websocket service
func (ws *websocketService) Name() string {
	return "websocket"
}

// AcceptConnection accept new connection
func (ws *websocketService) AcceptConnection(w http.ResponseWriter, r *http.Request, responseHeader http.Header) error {
	ws.logger.Debugf("new connection, peer: %s", r.RemoteAddr)

	conn, err := ws.Upgrade(w, r, responseHeader)

	if err != nil {
		ws.logger.Errorf("updrading error: %s", err)
		return err
	}

	wsConn := &websocketConn{
		ctx:           r.Context(),
		conn:          conn,
		readDeadline:  ws.config.HTTP.Websocket.ReadDeadline,
		writeDeadline: ws.config.HTTP.Websocket.WriteDeadline,
	}

	if err = ws.handleConnection(wsConn); err != nil {
		ws.handleError(err)
		return err
	}

	return nil
}

// handleConnection handle new connection
func (ws *websocketService) handleConnection(conn *websocketConn) error {
	defer ws.disconnect(conn)

	for {
		message, err := conn.readMessage()

		if err != nil {
			return err
		}

		response, err := ws.handleMessage(message, conn)
		conn.writeMessage(response)

		if err != nil {
			return fmt.Errorf("fail to handle message: %w", err)
		}
	}
}

// handleError handle websocket procession error
func (ws *websocketService) handleError(err error) {
	switch {
	case err == nil:
		return

	case websocket.IsCloseError(err, websocket.CloseGoingAway):
		ws.logger.Debugf("connection closed normal: %s", err.Error())

	case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
		ws.logger.Debugf("connection closed abnormal: %s", err.Error())

	default:
		ws.logger.Errorf("connection handle error: %s", err)
	}
}

// handleMessage handle incomming message
func (ws *websocketService) handleMessage(msg shared.WebsocketMessage, conn *websocketConn) (resp shared.WebsocketMessage, err error) {
	resp.Type = msg.Type
	resp.Payload = shared.ResponseOK

	ws.logger.Debugf("new message: %s", msg)

	if !conn.user.IsAuthorized {
		if err = ws.authorize(msg, conn); err != nil {
			resp.SetError(err)
			return
		}
	}

	switch msg.Type {
	case shared.WsMessagePing, shared.WsMessageAuth:
		resp.Type = shared.WsMessagePong

	case shared.WsMessagePong:
		resp.Type = shared.WsMessagePing

	case shared.WsMessageParams:
		if err = ws.updateParams(msg, conn); err != nil {
			resp.SetError(err)
		}

	default:
		err = fmt.Errorf("invalid message type: %s", msg.Type)
		resp.SetError(err)
	}

	return
}

// updateParams update user params
func (ws *websocketService) updateParams(msg shared.WebsocketMessage, conn *websocketConn) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if err := msg.Decode(&conn.user.Params); err != nil {
		return fmt.Errorf("error decode params: %w", err)
	}

	return nil
}

// authorize new connection
func (ws *websocketService) authorize(msg shared.WebsocketMessage, conn *websocketConn) error {
	if msg.Type != shared.WsMessageAuth {
		return fmt.Errorf("bad message type, have %s, want: %s", msg.Type, shared.WsMessageAuth)
	}

	if token, ok := msg.Payload.(string); !ok {
		return fmt.Errorf("bad token data: %s", token)

	} else if user, err := ws.authService.LoginByToken(conn.ctx, token, ws.config); err != nil {
		return fmt.Errorf("authorization error: %s", err.Error())

	} else {
		user.IsAuthorized = true
		conn.user = user

		ws.addConnection(conn)
	}

	return nil
}

// disconnect remove and close websocket connection
func (ws *websocketService) disconnect(conn *websocketConn) {
	ws.removeConnection(conn)
	conn.close()
}

// addConnection add new connection to local store
func (ws *websocketService) addConnection(conn *websocketConn) {
	conns, _ := ws.connections.Get(conn.user.ID)
	conns = append(conns, conn)
	ws.connections.Set(conn.user.ID, conns)

	ws.logger.Debugf("new authorized connection: %s", conn.user.GetLabel())
}

// removeConnection remove connection from local store
func (ws *websocketService) removeConnection(conn *websocketConn) {
	if !conn.user.IsAuthorized {
		return
	}

	conns, ok := ws.connections.Get(conn.user.ID)

	if !ok {
		return
	}

	conns = slices.DeleteFunc(conns, func(c *websocketConn) bool {
		return c == conn
	})

	if len(conns) > 0 {
		ws.connections.Set(conn.user.ID, conns)
	} else {
		ws.connections.Delete(conn.user.ID)
	}

	ws.logger.Debugf("client disconnected: %s", conn.user.GetLabel())
}

func (ws *websocketService) SubscribeEvent(event event.Event) error {
	msg := shared.WebsocketMessage{
		Type:    shared.WsMessageEvent,
		Payload: event,
	}

	for _, conns := range ws.connections.Values() {
		ws.sendToConnections(msg, conns...)
	}

	return nil
}

// FindAndSend sinf users by callback function and send message to it
func (ws *websocketService) FindAndSend(msg shared.WebsocketMessage, cb func(u user.CurrentUser) bool) {
	conns := ws.findConnections(cb)
	ws.sendToConnections(msg, conns...)
}

// findConnections find connections by callback function
func (ws *websocketService) findConnections(cb func(u user.CurrentUser) bool) (result []*websocketConn) {
	for _, conns := range ws.connections.Values() {
		for _, conn := range conns {
			if cb(conn.user) {
				result = append(result, conn)
			}
		}
	}

	return
}

// SendToUsers send message to users by ids
func (ws *websocketService) SendToUsers(msg shared.WebsocketMessage, ids ...int64) {
	for _, id := range ids {
		if conns, ok := ws.connections.Get(id); ok {
			ws.sendToConnections(msg, conns...)
		}
	}
}

// sendToConnections send message to connections
func (ws *websocketService) sendToConnections(msg shared.WebsocketMessage, conns ...*websocketConn) {
	for _, conn := range conns {
		conn.writeMessage(msg)
	}
}

// writeMessage send message to connection
func (c *websocketConn) writeMessage(msg shared.WebsocketMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.conn.SetWriteDeadline(time.Now().Add(c.writeDeadline))

	return c.conn.WriteJSON(msg)
}

// readMessage read new message from connection
func (c *websocketConn) readMessage() (msg shared.WebsocketMessage, err error) {
	c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))
	err = c.conn.ReadJSON(&msg)

	return
}

// close connection
func (c *websocketConn) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.Close()
}
