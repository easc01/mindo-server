package handlers

import (
	"database/sql"
	"fmt"
	HttpStatus "net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ishantSikdar/mindo-server/internal/constants"
	"github.com/ishantSikdar/mindo-server/internal/middleware"
	"github.com/ishantSikdar/mindo-server/internal/services"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/utils"
	"github.com/ishantSikdar/mindo-server/pkg/utils/http"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	userRg := rg.Group(constants.User, middleware.AuthMiddleware())

	{
		userRg.GET(utils.IdParam, getUserByID)
	}

}

func getUserByID(c *gin.Context) {
	paramId := c.Param("id")

	parsedId, parseErr := uuid.Parse(paramId)
	if parseErr != nil {
		http.NewErrorResponse(
			HttpStatus.StatusBadRequest,
			constants.InvalidUserID,
			parseErr,
		).Send(c)
		return
	}

	user, userErr := services.GetAppUserByUserID(parsedId)

	if userErr != nil {
		if userErr == sql.ErrNoRows {
			logger.Log.Errorf("user %s not found, %s", parsedId, userErr)
			http.NewErrorResponse(
				HttpStatus.StatusNotFound,
				fmt.Sprintf("User of %s ID not found", parsedId),
				userErr,
			).Send(c)
			return
		}

		logger.Log.Errorf("failed to get user %s userID: %s", userErr, parsedId)
		http.NewErrorResponse(
			HttpStatus.StatusInternalServerError,
			constants.SomethingWentWrong,
			userErr,
		).Send(c)
		return
	}

	http.NewResponse(
		HttpStatus.StatusFound,
		user,
	).Send(c)
}
