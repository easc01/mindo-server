package userhandler

import (
	"net/http"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	userservice "github.com/easc01/mindo-server/internal/services/user_service"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	httputil "github.com/easc01/mindo-server/pkg/utils/http_util"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterAdminUserRoutes(rg *gin.RouterGroup) {
	adminProtectedRg := rg.Group(route.Admin, middleware.RequireRole(models.UserTypeAdminUser))

	{
		adminProtectedRg.GET(constant.IdParam, getAdminUserByID)
	}

}

func getAdminUserByID(c *gin.Context) {
	paramId := c.Param("id")

	parsedId, parseErr := uuid.Parse(paramId)
	if parseErr != nil {
		httputil.NewErrorResponse(
			http.StatusBadRequest,
			message.InvalidUserID,
			parseErr.Error(),
		).Send(c)
		return
	}

	user, statusCode, userErr := userservice.GetAdminUserByUserID(parsedId)

	if userErr != nil {
		logger.Log.Errorf("failed to get admin %s id: %s", userErr, parsedId)
		httputil.NewErrorResponse(
			statusCode,
			message.SomethingWentWrong,
			userErr.Error(),
		).Send(c)
		return
	}

	httputil.NewResponse(
		http.StatusFound,
		user,
	).Send(c)
}
