package authservice

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/easc01/mindo-server/internal/config"
	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func IssueAuthTokens(c *gin.Context, userId uuid.UUID, role models.UserType) (string, error) {

	accessToken, atErr := CreateAccessToken(
		uuid.New().String(),
		userId.String(),
		role,
	)

	if atErr != nil {
		logger.Log.Errorf(
			"failed to create access token for userId: %s, %s",
			userId,
			atErr.Error(),
		)
		return constant.Blank, atErr
	}

	refreshToken, rtErr := CreateRefreshTokenByUserId(userId, role)
	if rtErr != nil {
		logger.Log.Errorf(
			"failed to create refresh token for userId: %s, %s",
			userId,
			rtErr.Error(),
		)
		return constant.Blank, rtErr
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constant.RefreshToken,
		Value:    refreshToken.RefreshToken.String(),
		HttpOnly: true,
		Secure:   config.GetConfig().Env == config.Production,
		Path:     route.GetRefreshRoute(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(constant.Month),
	})

	return accessToken, nil
}

func RefreshTokenService(c *gin.Context, refreshToken string) (dto.TokenDTO, int, error) {
	uuidRefreshToken, _ := uuid.Parse(refreshToken)

	// find refresh token
	userToken, utErr := db.Queries.GetUserTokenByRefreshToken(c, uuidRefreshToken)
	if utErr != nil {
		if errors.Is(utErr, sql.ErrNoRows) {
			logger.Log.Errorf("user token %s not found, %s", uuidRefreshToken, utErr.Error())
			return dto.TokenDTO{}, http.StatusUnauthorized, fmt.Errorf(message.SignInAgain)
		}
	}

	// check if expired
	if userToken.ExpiresAt.Before(time.Now()) {
		logger.Log.Errorf("refresh token is expired")
		return dto.TokenDTO{}, http.StatusUnauthorized, fmt.Errorf(message.SignInAgain)
	}

	// create tokens
	accessToken, tokenErr := IssueAuthTokens(c, userToken.UserID, userToken.Role)
	if tokenErr != nil {
		logger.Log.Errorf(
			"failed to issue auth tokens for user id: %s, %s",
			userToken.UserID.String(),
			tokenErr.Error(),
		)
		return dto.TokenDTO{}, http.StatusUnauthorized, fmt.Errorf(message.SignInAgain)
	}

	return dto.TokenDTO{
		AccessToken: accessToken,
	}, http.StatusAccepted, nil
}
