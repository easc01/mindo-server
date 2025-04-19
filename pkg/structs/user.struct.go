package structs

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/ishantSikdar/mindo-server/internal/models"
)

type GoogleLoginRequest struct {
	IDToken     string `json:"idToken"`
	AccessToken string `json:"accessToken"`
}

type AppUserDataDTO struct {
	UserID            uuid.UUID
	UserType          models.UserType
	Username          sql.NullString
	ProfilePictureUrl sql.NullString
	Bio               sql.NullString
	Name              sql.NullString
	Mobile            sql.NullString
	Email             sql.NullString
	LastLoginAt       sql.NullTime
	UpdatedAt         sql.NullTime
	CreatedAt         sql.NullTime
	UpdatedBy         uuid.NullUUID
}

type NewAppUserParams struct {
	Name     string
	Username string
	Email    string
	Mobile   string
}
