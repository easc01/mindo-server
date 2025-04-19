package handlers

import (
	HttpStatus "net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/internal/constants"
	"github.com/ishantSikdar/mindo-server/internal/services"
	"github.com/ishantSikdar/mindo-server/pkg/structs"
	"github.com/ishantSikdar/mindo-server/pkg/utils/http"
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
		http.NewErrorResponse(
			HttpStatus.StatusBadRequest,
			constants.InvalidRequestBody,
			err,
		).Send(c)
		return
	}

	user, userErr := services.GoogleAuthService(c, req)

	if userErr != nil {
		http.NewErrorResponse(
			HttpStatus.StatusInternalServerError,
			constants.SomethingWentWrong,
			userErr,
		).Send(c)
		return
	}

	http.NewResponse(
		HttpStatus.StatusCreated,
		user,
	).Send(c)
}
