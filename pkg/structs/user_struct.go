package structs

import (
	"time"

	"github.com/google/uuid"
	"github.com/ishantSikdar/mindo-server/internal/models"
)

type GoogleLoginRequest struct {
	IDToken     string `json:"idToken"`
	AccessToken string `json:"accessToken"`
}

type AppUserDataDTO struct {
	UserID            uuid.UUID       `json:"userId"`
	UserType          models.UserType `json:"userType"`
	Username          string          `json:"username"`
	ProfilePictureUrl string          `json:"profilePictureUrl"`
	OauthClientID     string          `json:"oauthClientId"`
	Bio               string          `json:"bio"`
	Name              string          `json:"name"`
	Mobile            string          `json:"mobile"`
	Email             string          `json:"email"`
	LastLoginAt       time.Time       `json:"lastLoginAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedBy         uuid.UUID       `json:"updatedBy"`
}

type NewAppUserParams struct {
	Name          string
	Username      string
	Email         string
	Mobile        string
	OauthClientID string
}
