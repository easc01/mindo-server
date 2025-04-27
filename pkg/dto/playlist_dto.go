package dto

import "time"

type CreatePlaylistRequest struct {
	Name         string   `json:"name"         binding:"required"`
	Description  string   `json:"description"  binding:"required"`
	DomainName   string   `json:"domainName"   binding:"required"`
	ThumbnailURL string   `json:"thumbnailUrl" binding:"required"`
	Topics       []string `json:"topics"       binding:"required,dive"`
}

type PlaylistDetailsDTO struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	InterestID   string    `json:"interestId"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	Views        int       `json:"views"`
	Code         string    `json:"code"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	UpdatedBy    string    `json:"updatedBy"`
	Topics       []string  `json:"topics"`
}

type PlaylistPreviewDTO struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	InterestID   string    `json:"interestId"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	Views        int       `json:"views"`
	Code         string    `json:"code"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	UpdatedBy    string    `json:"updatedBy"`
	TopicsCount  int       `json:"topicsCount"`
}
