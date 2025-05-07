package dto

import (
	"time"

	"github.com/gorilla/websocket"
)

type ChatClient struct {
	Conn    *websocket.Conn
	AppUser *AppUserDataDTO
}

type ChatMessage struct {
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
