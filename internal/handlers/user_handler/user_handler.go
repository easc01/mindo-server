package userhandler

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ishantSikdar/mindo-server/internal/middleware"
	userservice "github.com/ishantSikdar/mindo-server/internal/services/user_service"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/utils/constant"
	"github.com/ishantSikdar/mindo-server/pkg/utils/http"
	"github.com/ishantSikdar/mindo-server/pkg/utils/message"
	"github.com/ishantSikdar/mindo-server/pkg/utils/route"
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
			parseErr,
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
				userErr,
			).Send(c)
			return
		}

		logger.Log.Errorf("failed to get user %s userID: %s", userErr, parsedId)
		http.NewErrorResponse(
			http.StatusInternalServerError,
			message.SomethingWentWrong,
			userErr,
		).Send(c)
		return
	}

	http.NewResponse(
		http.StatusFound,
		user,
	).Send(c)
}
