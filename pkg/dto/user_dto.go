package dto

import (
	"time"

	"github.com/easc01/mindo-server/internal/models"
	"github.com/google/uuid"
)

type AppUserDataDTO struct {
	AccessToken       string          `json:"accessToken,omitempty"`
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

type AdminUserDataDTO struct {
	AccessToken string          `json:"accessToken,omitempty"`
	UserID      uuid.UUID       `json:"userId"`
	UserType    models.UserType `json:"userType"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	LastLoginAt time.Time       `json:"lastLoginAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedBy   uuid.UUID       `json:"updatedBy"`
}

type NewAppUserParams struct {
	Name          string
	Username      string
	Email         string
	Mobile        string
	OauthClientID string
}

type NewAdminUserParams struct {
	Name     string `json:"name"     binding:"required,min=8"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type AdminSignInParams struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type TokenDTO struct {
	AccessToken string `json:"accessToken"`
}
