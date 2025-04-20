package authhandler

import (
	authservice "github.com/easc01/mindo-server/internal/services/auth_service"
	userservice "github.com/easc01/mindo-server/internal/services/user_service"
	"github.com/easc01/mindo-server/pkg/dto"
	httputil "github.com/easc01/mindo-server/pkg/utils/http_util"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
)

func RegisterAuth(rg *gin.RouterGroup) {
	authRg := rg.Group(route.Auth)

	{
		authRg.POST(route.Google, googleAuthHandler)
		authRg.POST(route.Admin, adminSignUpHandler)
		authRg.POST(route.Admin+"/sign-in", adminSignInHandler)
	}
}

func googleAuthHandler(c *gin.Context) {
	req, ok := httputil.GetRequestBody[dto.GoogleLoginRequest](c)
	if !ok {
		return
	}

	user, statusCode, userErr := authservice.GoogleAuthService(c, &req)

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

func adminSignUpHandler(c *gin.Context) {
	req, ok := httputil.GetRequestBody[dto.NewAdminUserParams](c)
	if !ok {
		return
	}

	user, userErr := userservice.CreateNewAdminUser(&req)

	if userErr != nil {
		httputil.NewErrorResponse(
			httputil.StatusInternalServerError,
			message.SomethingWentWrong,
			userErr.Error(),
		).Send(c)
		return
	}

	httputil.NewResponse(
		httputil.StatusCreated,
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
