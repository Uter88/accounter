package protocols

import (
	"accounter/config"
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

// Websocket messages type
type wsMessageType string

const (
	wsMessageAuth wsMessageType = "authorize"
	wsMessagePing wsMessageType = "ping"
	wsMessagePong wsMessageType = "pong"
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

	// Websocket message
	websocketMessage struct {
		Type    wsMessageType `json:"type"`
		Payload any           `json:"payload"`
		Error   string        `json:"error"`
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

// AcceptConnection accept new connection
func (ws *websocketService) AcceptConnection(w http.ResponseWriter, r *http.Request, responseHeader http.Header) error {
	conn, err := ws.Upgrade(w, r, responseHeader)

	if err != nil {
		ws.logger.Errorf("WS: updrading error: %s", err.Error())
		return err
	}

	wsConn := &websocketConn{
		ctx:           r.Context(),
		conn:          conn,
		readDeadline:  ws.config.HTTP.Websocket.ReadDeadline,
		writeDeadline: ws.config.HTTP.Websocket.WriteDeadline,
	}

	err = ws.handleConnection(wsConn)

	if err = ws.handleConnection(wsConn); err != nil {
		ws.logger.Errorf("WS: handle connection error: %s", err.Error())
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
			return fmt.Errorf("error to read message: %s", err.Error())
		}

		if !conn.user.IsAuthorized {
			if err = ws.authorize(message, conn); err != nil {
				return err
			}
		}

		response := ws.handleMessage(message)

		if err = conn.writeMessage(response); err != nil {
			return fmt.Errorf("error of write message: %s", err.Error())
		}
	}
}

// handleMessage handle incomming message
func (ws *websocketService) handleMessage(msg websocketMessage) (resp websocketMessage) {
	resp.Type = msg.Type

	switch msg.Type {
	case wsMessagePing:
		resp.Type = wsMessagePong

	case wsMessagePong:
		resp.Type = wsMessagePing

	case wsMessageAuth:
		resp.Error = "user is already authorized"

	default:
		resp.Error = fmt.Sprintf("invalid message type: %s", msg.Type)
	}

	return
}

// authorize new connection
func (ws *websocketService) authorize(msg websocketMessage, conn *websocketConn) error {
	if msg.Type != wsMessageAuth {
		return fmt.Errorf("bad message type, have %s, want: %s", msg.Type, wsMessageAuth)
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

	ws.logger.Infof("WS: new connection: %s", conn.user.GetLabel())
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

	ws.logger.Infof("WS: client disconnected: %s", conn.user.GetLabel())
}

// writeMessage send message to connection
func (c *websocketConn) writeMessage(msg websocketMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.conn.SetWriteDeadline(time.Now().Add(c.writeDeadline))

	return c.conn.WriteJSON(msg)
}

// readMessage read new message from connection
func (c *websocketConn) readMessage() (msg websocketMessage, err error) {
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
