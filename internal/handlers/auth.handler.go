package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/internal/constants"
	"github.com/ishantSikdar/mindo-server/internal/services"
	"github.com/ishantSikdar/mindo-server/pkg/structs"
)

func RegisterAuth(rg *gin.RouterGroup) {
	authRg := rg.Group(constants.Auth)

	{
		authRg.POST(constants.Google, googleAuthHandler)
	}
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

	c.JSON(http.StatusCreated, user)
}
