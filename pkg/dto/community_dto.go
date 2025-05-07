package dto

import (
	"time"

	"github.com/google/uuid"
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

type CommunitiesDTO struct {
	ID           uuid.UUID     `json:"id"`
	Title        string        `json:"title"`
	About        string        `json:"about"`
	ThumbnailUrl string        `json:"thumbnailUrl"`
	LogoUrl      string        `json:"logoUrl"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedBy    uuid.NullUUID `json:"updatedBy"`
}
