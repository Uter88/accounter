package shared

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Websocket messages type
type WsMessageType = string

const (
	WsMessageAuth   WsMessageType = "authorize"
	WsMessagePing   WsMessageType = "ping"
	WsMessagePong   WsMessageType = "pong"
	WsMessageParams WsMessageType = "params"
	WsMessageEvent  WsMessageType = "event"
)

// Websocket message
type WebsocketMessage struct {
	Type    WsMessageType `json:"type"`
	Payload any           `json:"payload,omitempty"`
	Error   string        `json:"error,omitempty"`
}

func (m *WebsocketMessage) SetError(err error) {
	if err != nil {
		m.Error = err.Error()
		m.Payload = ResponseError
	}
}

// String representation of WebsocketMessage
func (m WebsocketMessage) String() string {
	return fmt.Sprintf("Type: %s, payload: %v", m.Type, m.Payload)
}

func (m WebsocketMessage) ToString() string {
	data, _ := json.Marshal(m)

	return string(data)
}

func (m WebsocketMessage) Decode(dest any) error {
	var content []byte

	switch tp := m.Payload.(type) {
	case []byte:
		content = tp
	case string:
		content = []byte(tp)
	case map[string]any:
		content, _ = json.Marshal(m.Payload)
	default:
		return errors.New("invalid payload type")
	}

	return json.Unmarshal(content, dest)
}
