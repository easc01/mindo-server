package authhandler

import (
	authservice "github.com/easc01/mindo-server/internal/services/auth_service"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/utils/http"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
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
