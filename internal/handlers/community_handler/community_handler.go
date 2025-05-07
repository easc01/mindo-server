package communityhandler

import (
	"net/http"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	communityservice "github.com/easc01/mindo-server/internal/services/community_service"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/message"
	networkutil "github.com/easc01/mindo-server/pkg/utils/network_util"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterCommunity(rg *gin.RouterGroup) {
	communityRg := rg.Group(route.Communities)

	{
		communityRg.POST(
			constant.Blank,
			middleware.RequireRole(models.UserTypeAppUser),
			createCommunity,
		)
		communityRg.POST(
			"/join/:communityId",
			middleware.RequireRole(models.UserTypeAppUser),
			joinCommunity,
		)
	}
}

func createCommunity(c *gin.Context) {
	req, ok := networkutil.GetRequestBody[dto.CreateCommunityDTO](c)
	if !ok {
		return
	}

	community, statusCode, err := communityservice.CreateNewCommunity(c, &req)

	if err != nil {
		networkutil.NewErrorResponse(
			statusCode,
			message.SomethingWentWrong,
			err.Error(),
		).Send(c)
		return
	}

	networkutil.NewResponse(
		statusCode,
		community,
	).Send(c)
}

func joinCommunity(c *gin.Context) {
	communityId := c.Param("communityId")

	parsedCommId, err := uuid.Parse(communityId)
	if err != nil {
		networkutil.NewErrorResponse(
			http.StatusBadRequest,
			"invalid community id",
			err.Error(),
		).Send(c)
		return
	}

	statusCode, err := communityservice.JoinExistingCommunity(c, parsedCommId)

	if err != nil {
		networkutil.NewErrorResponse(
			statusCode,
			message.SomethingWentWrong,
			err.Error(),
		).Send(c)
		return
	}

	networkutil.NewResponse(
		statusCode,
		"community joined",
	).Send(c)
}
