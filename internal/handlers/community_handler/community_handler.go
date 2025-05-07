package communityhandler

import (
	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	communityservice "github.com/easc01/mindo-server/internal/services/community_service"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/message"
	networkutil "github.com/easc01/mindo-server/pkg/utils/network_util"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
)

func RegisterCommunity(rg *gin.RouterGroup) {
	communityRg := rg.Group(route.Communities)

	{
		communityRg.POST(
			constant.Blank,
			middleware.RequireRole(models.UserTypeAppUser, models.UserTypeAdminUser),
			createCommunity,
		)
		communityRg.POST(
			"/join/:communityId",
			middleware.RequireRole(models.UserTypeAppUser, models.UserTypeAdminUser),
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

}
