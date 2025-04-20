package authservice

import (
	"net/http"

	"github.com/easc01/mindo-server/internal/config"
	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func IssueAuthTokens(c *gin.Context, userId uuid.UUID) (string, error) {

	accessToken, atErr := CreateAccessToken(
		uuid.New().String(),
		userId.String(),
		models.UserTypeAppUser,
	)

	if atErr != nil {
		logger.Log.Errorf(
			"failed to create access token for userId: %s, %s",
			userId,
			atErr.Error(),
		)
		return constant.Blank, atErr
	}

	refreshToken, rtErr := CreateRefreshTokenByUserId(userId)
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
		Value:    refreshToken.RefreshToken,
		HttpOnly: true,
		Secure:   config.GetConfig().Env == config.Production,
		Path:     route.GetRefreshRoute(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(constant.Month),
	})

	return accessToken, nil
}
