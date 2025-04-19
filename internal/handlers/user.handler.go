package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ishantSikdar/mindo-server/pkg/db"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/utils"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	userRg := rg.Group("/user")

	{
		userRg.GET(utils.IdParam, getUserByID)
	}

}

func getUserByID(c *gin.Context) {
	paramId := c.Param("id")

	parsedId, userErr := uuid.Parse(paramId)
	if userErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, userErr := db.Queries.GetAppUserByUserID(context.Background(), parsedId)

	if userErr != nil {
		if userErr == sql.ErrNoRows {
			logger.Log.Errorf("user %s not found, %s", parsedId, userErr)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		logger.Log.Errorf("failed to get user %s userID: %s", userErr, parsedId)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong while fetching user",
		})
		return
	}


	c.JSON(http.StatusFound, user)
}
