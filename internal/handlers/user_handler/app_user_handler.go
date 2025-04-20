package userhandler

import (
	"github.com/easc01/mindo-server/internal/middleware"
	userservice "github.com/easc01/mindo-server/internal/services/user_service"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/http"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterAppUserRoutes(rg *gin.RouterGroup) {
	userRg := rg.Group(route.User, middleware.AuthMiddleware())

	{
		userRg.GET(constant.IdParam, getAppUserByID)
	}

}

func getAppUserByID(c *gin.Context) {
	paramId := c.Param("id")

	parsedId, parseErr := uuid.Parse(paramId)
	if parseErr != nil {
		http.NewErrorResponse(
			http.StatusBadRequest,
			message.InvalidUserID,
			parseErr.Error(),
		).Send(c)
		return
	}

	user, statusCode, userErr := userservice.GetAppUserByUserID(parsedId)

	if userErr != nil {
		logger.Log.Errorf("failed to get user %s userID: %s", userErr, parsedId)
		http.NewErrorResponse(
			statusCode,
			message.SomethingWentWrong,
			userErr.Error(),
		).Send(c)
		return
	}

	http.NewResponse(
		http.StatusFound,
		user,
	).Send(c)
}
