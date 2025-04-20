package authservice

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/easc01/mindo-server/internal/config"
	"github.com/easc01/mindo-server/internal/models"
	jwtservice "github.com/easc01/mindo-server/internal/services/jwt_service"
	userservice "github.com/easc01/mindo-server/internal/services/user_service"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

func GoogleAuthService(
	c *gin.Context,
	googleReq *dto.GoogleLoginRequest,
) (dto.AppUserDataDTO, int, error) {

	payload, payloadErr := idtoken.Validate(c, googleReq.IDToken, config.GetConfig().GoogleClientId)
	if payloadErr != nil {
		logger.Log.Errorf(
			"invalid app user token, %s, for token %s",
			payloadErr,
			googleReq.IDToken,
		)
		return dto.AppUserDataDTO{}, http.StatusBadRequest, payloadErr
	}

	name, _ := payload.Claims["name"].(string)
	email, _ := payload.Claims["email"].(string)

	appUserParams := dto.NewAppUserParams{
		Name:          name,
		Email:         email,
		OauthClientID: payload.Subject,
		Username:      util.GenerateUsername(),
		Mobile:        constant.Blank,
	}

	// Check if appUser exists by oauthclientId and update last login
	appUser, appUserErr := db.Queries.UpdateAppUserLastLoginAtByOAuthClientID(
		c,
		util.GetSQLNullString(appUserParams.OauthClientID),
	)

	if appUserErr != nil {
		if errors.Is(appUserErr, sql.ErrNoRows) {
			// Create new user
			newAppUser, newAppUserErr := userservice.CreateNewAppUser(appUserParams)
			if newAppUserErr != nil {
				logger.Log.Errorf(
					"failed to create app user %s for email %s oauth client id %s",
					newAppUserErr,
					appUserParams.Email,
					appUserParams.OauthClientID,
				)
				return dto.AppUserDataDTO{}, http.StatusInternalServerError, newAppUserErr
			}

			logger.Log.Infof("new app user created %s", newAppUser.UserID)
			return newAppUser, http.StatusCreated, nil
		}

		// Log unexpected DB errors
		logger.Log.Errorf(
			"failed to update last login %s for oauth client id %s",
			appUserErr,
			appUserParams.OauthClientID,
		)
		return dto.AppUserDataDTO{}, http.StatusInternalServerError, appUserErr
	}

	// create tokens
	accessToken, atErr := jwtservice.CreateAccessToken(
		uuid.New().String(),
		appUser.UserID.String(),
		models.UserTypeAppUser,
	)

	if atErr != nil {
		logger.Log.Errorf(
			"failed to create access token for userId: %s, %s",
			appUser.UserID,
			atErr.Error(),
		)
		return dto.AppUserDataDTO{}, http.StatusInternalServerError, atErr
	}

	refreshToken, rtErr := jwtservice.CreateRefreshTokenByUserId(appUser.UserID)
	if rtErr != nil {
		logger.Log.Errorf(
			"failed to create refresh token for userId: %s, %s",
			appUser.UserID,
			rtErr.Error(),
		)
		return dto.AppUserDataDTO{}, http.StatusInternalServerError, rtErr
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constant.RefreshToken,
		Value:    refreshToken.RefreshToken,
		HttpOnly: true,
		Secure:   config.GetConfig().Env == config.Production,
		Path:     route.GetRefreshRoute(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(constant.Month),
	})

	return dto.AppUserDataDTO{
		AccessToken:       accessToken,
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
	}, http.StatusFound, nil
}
