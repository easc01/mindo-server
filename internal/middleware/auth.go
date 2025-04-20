package middleware

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/internal/config"
	"github.com/ishantSikdar/mindo-server/internal/models"
	"github.com/ishantSikdar/mindo-server/pkg/db"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/dto"
	"github.com/ishantSikdar/mindo-server/pkg/utils/constant"
	"github.com/ishantSikdar/mindo-server/pkg/utils/http"
	"github.com/ishantSikdar/mindo-server/pkg/utils/message"
	"github.com/ishantSikdar/mindo-server/pkg/utils/util"
	"google.golang.org/api/idtoken"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(constant.Authorization)

		if token == "" {
			http.NewErrorResponse(
				http.StatusUnauthorized,
				message.AuthHeaderRequired,
				message.ProvideAuthHeader,
			).Send(c)
			c.Abort()
			return
		}

		payload, payloadErr := idtoken.Validate(c, token, config.GetConfig().GoogleClientId)
		if payloadErr != nil {
			logger.Log.Errorf("invalid auth token %s", payloadErr)
			http.NewErrorResponse(
				http.StatusUnauthorized,
				fmt.Sprintf("Invalid auth token, %s", payloadErr.Error()),
				payloadErr,
			).Send(c)
			c.Abort()
			return
		}

		appUser, appUserErr := db.Queries.GetAppUserByClientOAuthID(
			c.Request.Context(),
			util.GetSQLNullString(payload.Subject),
		)
		if appUserErr != nil {
			if errors.Is(appUserErr, sql.ErrNoRows) {
				http.NewErrorResponse(
					http.StatusNotFound,
					message.UserNotFound,
					appUserErr,
				).Send(c)
			} else {
				http.NewErrorResponse(
					http.StatusInternalServerError,
					message.SomethingWentWrong,
					appUserErr,
				).Send(c)
			}
			logger.Log.Errorf(
				"failed to get user of auth client id: %s, %s",
				payload.Subject,
				appUserErr,
			)
			c.Abort()
			return
		}

		appUserContext := dto.AppUserDataDTO{
			UserID:            appUser.UserID,
			Username:          appUser.Username.String,
			ProfilePictureUrl: appUser.ProfilePictureUrl.String,
			OauthClientID:     appUser.OauthClientID.String,
			Bio:               appUser.Bio.String,
			Name:              appUser.Name.String,
			Mobile:            appUser.Mobile.String,
			Email:             appUser.Email.String,
			LastLoginAt:       appUser.LastLoginAt.Time,
			UpdatedAt:         appUser.UpdatedAt.Time,
			CreatedAt:         appUser.CreatedAt.Time,
			UpdatedBy:         appUser.UpdatedBy.UUID,
			UserType:          models.UserTypeAppUser,
		}

		c.Set(string(UserContextKey), appUserContext)
		c.Next()
	}
}

func GetUser(ctx *gin.Context) (dto.AppUserDataDTO, bool) {
	user, ok := ctx.Get(string(UserContextKey))
	if !ok {
		return dto.AppUserDataDTO{}, false
	}
	appUser, ok := user.(dto.AppUserDataDTO)
	return appUser, ok
}
