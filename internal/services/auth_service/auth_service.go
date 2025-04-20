package authservice

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/easc01/mindo-server/internal/config"
	"github.com/easc01/mindo-server/internal/models"
	userservice "github.com/easc01/mindo-server/internal/services/user_service"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
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
	accessToken, tokenErr := IssueAuthTokens(c, appUser.UserID)
	if tokenErr != nil {
		logger.Log.Errorf(
			"failed to issue auth tokens for user id: %s, %s",
			appUser.UserID.String(),
			tokenErr.Error(),
		)
		return dto.AppUserDataDTO{}, http.StatusInternalServerError, tokenErr
	}

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

func RefreshTokenService(c *gin.Context, refreshToken string) (dto.TokenResponse, int, error) {
	// find refresh token
	userToken, utErr := db.Queries.GetUserTokenByRefreshToken(c, refreshToken)
	if utErr != nil {
		if errors.Is(utErr, sql.ErrNoRows) {
			logger.Log.Errorf("user token not found, %s", utErr.Error())
			return dto.TokenResponse{}, http.StatusUnauthorized, fmt.Errorf(message.SignInAgain)
		}
	}

	// check if expired
	if userToken.ExpiresAt.Before(time.Now()) {
		logger.Log.Errorf("refresh token is expired")
		return dto.TokenResponse{}, http.StatusUnauthorized, fmt.Errorf(message.SignInAgain)
	}

	// create tokens
	accessToken, tokenErr := IssueAuthTokens(c, userToken.UserID)
	if tokenErr != nil {
		logger.Log.Errorf(
			"failed to issue auth tokens for user id: %s, %s",
			userToken.UserID.String(),
			tokenErr.Error(),
		)
		return dto.TokenResponse{}, http.StatusUnauthorized, fmt.Errorf(message.SignInAgain)
	}

	return dto.TokenResponse{
		AccessToken: accessToken,
	}, http.StatusAccepted, nil
}
