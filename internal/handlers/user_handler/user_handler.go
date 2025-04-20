package userhandler

import (
	"database/sql"
	"fmt"

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

func RegisterUserRoutes(rg *gin.RouterGroup) {
	userRg := rg.Group(route.User, middleware.AuthMiddleware())

	{
		userRg.GET(constant.IdParam, getUserByID)
	}

}

func getUserByID(c *gin.Context) {
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

	user, userErr := userservice.GetAppUserByUserID(parsedId)

	if userErr != nil {
		if userErr == sql.ErrNoRows {
			logger.Log.Errorf("user %s not found, %s", parsedId, userErr)
			http.NewErrorResponse(
				http.StatusNotFound,
				fmt.Sprintf("User of %s ID not found", parsedId),
				userErr.Error(),
			).Send(c)
			return
		}

		logger.Log.Errorf("failed to get user %s userID: %s", userErr, parsedId)
		http.NewErrorResponse(
			http.StatusInternalServerError,
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
