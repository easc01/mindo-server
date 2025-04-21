package userhandler

import (
	"net/http"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	userservice "github.com/easc01/mindo-server/internal/services/user_service"
	"github.com/easc01/mindo-server/pkg/dto"
	httputil "github.com/easc01/mindo-server/pkg/utils/http_util"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
)

func RegisterAdminUserRoutes(rg *gin.RouterGroup) {
	adminRg := rg.Group(route.Admin, middleware.RequireRole(models.UserTypeAdminUser))

	{
		adminRg.POST(route.Admin, adminSignUpHandler)
		adminRg.POST(route.Admin+"/sign-in", adminSignInHandler)
	}

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
