package networkutil

import "github.com/gorilla/websocket"

type WSMessage struct {
	Code    int `json:"code,omitempty"`
	Message string `json:"message"`
}

func WSError(
	code int,
	message string,
) *WSMessage {
	return &WSMessage{
		Code:    code,
		Message: message,
	}
}

func (ws *WSMessage) Send(
	conn *websocket.Conn,
) {
	conn.WriteJSON(ws)
	conn.Close()
}
