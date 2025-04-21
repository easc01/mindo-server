package authhandler

import (
	"net/http"

	authservice "github.com/easc01/mindo-server/internal/services/auth_service"
	userservice "github.com/easc01/mindo-server/internal/services/user_service"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	httputil "github.com/easc01/mindo-server/pkg/utils/http_util"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
)

func RegisterAuth(rg *gin.RouterGroup) {
	authRg := rg.Group(route.Auth)
	adminAuthRg := rg.Group(route.Auth + route.Admin)

	{
		authRg.POST(route.Google, googleAuthHandler)
		authRg.POST(route.Refresh, refreshTokenHandler)
	}

	{
		adminAuthRg.POST(route.SignUp, adminSignUpHandler)
		adminAuthRg.POST(route.SignIn, adminSignInHandler)
	}
}

func googleAuthHandler(c *gin.Context) {
	req, ok := httputil.GetRequestBody[dto.GoogleLoginRequest](c)
	if !ok {
		return
	}

	user, statusCode, userErr := userservice.GoogleAuthService(c, &req)

	if userErr != nil {
		httputil.NewErrorResponse(
			statusCode,
			message.SomethingWentWrong,
			userErr.Error(),
		).Send(c)
		return
	}

	httputil.NewResponse(
		statusCode,
		user,
	).Send(c)
}

func refreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie(constant.RefreshToken)
	if err != nil {
		logger.Log.Errorf("%s cookie not found", constant.RefreshToken)
		httputil.NewErrorResponse(
			http.StatusUnauthorized,
			message.SignInAgain,
			nil,
		)
	}

	token, statusCode, err := authservice.RefreshTokenService(c, refreshToken)
	if err != nil {
		logger.Log.Errorf("failed to generate access and refresh tokens, %s", err.Error())
		httputil.NewErrorResponse(
			statusCode,
			err.Error(),
			constant.Blank,
		).Send(c)
	}

	httputil.NewResponse(
		statusCode,
		token,
	).Send(c)
}

func adminSignUpHandler(c *gin.Context) {
	req, ok := httputil.GetRequestBody[dto.NewAdminUserParams](c)
	if !ok {
		return
	}

	user, userErr := userservice.CreateNewAdminUser(&req)

	if userErr != nil {
		httputil.NewErrorResponse(
			http.StatusInternalServerError,
			message.SomethingWentWrong,
			userErr.Error(),
		).Send(c)
		return
	}

	httputil.NewResponse(
		http.StatusCreated,
		user,
	).Send(c)
}

func adminSignInHandler(c *gin.Context) {
	req, ok := httputil.GetRequestBody[dto.AdminSignInParams](c)
	if !ok {
		return
	}

	user, statusCode, userErr := userservice.AdminSignIn(c, &req)

	if userErr != nil {
		httputil.NewErrorResponse(
			statusCode,
			message.SomethingWentWrong,
			userErr.Error(),
		).Send(c)
		return
	}

	httputil.NewResponse(
		statusCode,
		user,
	).Send(c)
}
