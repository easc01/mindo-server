package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/internal/constants"
	"github.com/ishantSikdar/mindo-server/internal/services"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/structs"
	"github.com/ishantSikdar/mindo-server/pkg/utils"
)

func RegisterAuth(rg *gin.RouterGroup) {
	authRg := rg.Group(constants.Auth)

	{
		authRg.POST(constants.Google, googleAuthHandler)
	}
	logger.Log.Info("registered auth routes")
}

func googleAuthHandler(c *gin.Context) {
	var req structs.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, userErr := services.GoogleAuthService(c, req)

	if userErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	userParsed, userParsedErr := utils.ParseSQLResponse(user)
	if userParsedErr != nil {
		logger.Log.Error("failed to parsed user data", userParsedErr)
	}

	c.JSON(http.StatusOK, userParsed)
}
