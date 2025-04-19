package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/internal/config"
	"github.com/ishantSikdar/mindo-server/internal/models"
	"github.com/ishantSikdar/mindo-server/pkg/db"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/structs"
	"github.com/ishantSikdar/mindo-server/pkg/utils"
	"google.golang.org/api/idtoken"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(utils.Authorization)

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		payload, payloadErr := idtoken.Validate(c, token, config.GetConfig().GoogleClientId)
		if payloadErr != nil {
			logger.Log.Errorf("invalid auth token %s", payloadErr)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token, " + payloadErr.Error()})
			c.Abort()
			return
		}

		appUser, appUserErr := db.Queries.GetAppUserByClientOAuthID(
			c.Request.Context(),
			utils.GetSQLNullString(payload.Subject),
		)
		if appUserErr != nil {
			if errors.Is(appUserErr, sql.ErrNoRows) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			logger.Log.Errorf(
				"failed to get user of auth client id: %s, %s",
				payload.Subject,
				appUserErr,
			)
			c.Abort()
			return
		}

		appUserContext := structs.AppUserDataDTO{
			UserID:            appUser.UserID,
			Username:          appUser.Username,
			ProfilePictureUrl: appUser.ProfilePictureUrl,
			Bio:               appUser.Bio,
			OauthClientID:     appUser.OauthClientID,
			Name:              appUser.Name,
			Mobile:            appUser.Mobile,
			Email:             appUser.Email,
			LastLoginAt:       appUser.LastLoginAt,
			UpdatedAt:         appUser.UpdatedAt,
			CreatedAt:         appUser.CreatedAt,
			UpdatedBy:         appUser.UpdatedBy,
			UserType:          models.UserTypeAppUser,
		}

		c.Set(string(UserContextKey), appUserContext)
		c.Next()
	}
}

func GetUser(ctx *gin.Context) (structs.AppUserDataDTO, bool) {
	user, ok := ctx.Get(string(UserContextKey))
	if !ok {
		return structs.AppUserDataDTO{}, false
	}
	appUser, ok := user.(structs.AppUserDataDTO)
	return appUser, ok
}
