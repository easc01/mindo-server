package dto

import (
	"time"

	"github.com/easc01/mindo-server/internal/models"
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

type CommunityDTO struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	About        string    `json:"about"`
	ThumbnailUrl string    `json:"thumbnailUrl"`
	LogoUrl      string    `json:"logoUrl"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedBy    string    `json:"updatedBy"`
}

type CreateCommunityDTO struct {
	Title        string `json:"title"        binding:"required"`
	About        string `json:"about"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	LogoUrl      string `json:"logoUrl"`
}

type UserMessageDTO struct {
	MessageGroupID uuid.UUID    `json:"messageGroupId"`
	UserID         uuid.UUID    `json:"userId"`
	Name           string       `json:"name"`
	Username       string       `json:"username"`
	UserProfileUrl string       `json:"userProfilePic"`
	UserColor      models.Color `json:"userColor"`
	Messages       []MessageDTO `json:"messages"`
}

type MessageDTO struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type SocketMessageDTO struct {
	MessageGroupID uuid.UUID    `json:"messageGroupId"`
	MessageId      uuid.UUID    `json:"messageId"`
	UserID         uuid.UUID    `json:"userId"`
	Name           string       `json:"name"`
	Username       string       `json:"username"`
	UserProfileUrl string       `json:"userProfilePic"`
	UserColor      models.Color `json:"userColor"`
	CommunityID    uuid.UUID    `json:"communityId"`
	Content        string       `json:"content"`
	Timestamp      time.Time    `json:"timestamp"`
}
