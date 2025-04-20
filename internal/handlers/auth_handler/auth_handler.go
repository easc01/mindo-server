package authhandler

import (
	"github.com/gin-gonic/gin"
	authservice "github.com/ishantSikdar/mindo-server/internal/services/auth_service"
	"github.com/ishantSikdar/mindo-server/pkg/dto"
	"github.com/ishantSikdar/mindo-server/pkg/utils/http"
	"github.com/ishantSikdar/mindo-server/pkg/utils/message"
	"github.com/ishantSikdar/mindo-server/pkg/utils/route"
)

func RegisterAuth(rg *gin.RouterGroup) {
	authRg := rg.Group(route.Auth)

	{
		authRg.POST(route.Google, googleAuthHandler)
	}
}

func googleAuthHandler(c *gin.Context) {
	req, ok := http.GetRequestBody[dto.GoogleLoginRequest](c)

	if !ok {
		return
	}

	user, userErr := authservice.GoogleAuthService(c, &req)

	if userErr != nil {
		http.NewErrorResponse(
			http.StatusInternalServerError,
			message.SomethingWentWrong,
			userErr,
		).Send(c)
		return
	}

	http.NewResponse(
		http.StatusCreated,
		user,
	).Send(c)
}
