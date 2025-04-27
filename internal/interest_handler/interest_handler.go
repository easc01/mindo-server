package interesthandler

import (
	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	interestservice "github.com/easc01/mindo-server/internal/services/interest_service"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	httputil "github.com/easc01/mindo-server/pkg/utils/http_util"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
)

func RegisterInterest(rg *gin.RouterGroup) {
	intRg := rg.Group(route.Interest)

	{
		intRg.POST(
			constant.Blank,
			middleware.RequireRole(models.UserTypeAdminUser),
			upsertMasterInterestHandler,
		)
		intRg.GET(constant.Blank, getMasterInterestListHandler)
	}
}

func upsertMasterInterestHandler(c *gin.Context) {
	
}

func getMasterInterestListHandler(c *gin.Context) {
	interests, statusCode, intErr := interestservice.GetMasterInterestList(c)

	if intErr != nil {
		httputil.NewErrorResponse(statusCode, message.SomethingWentWrong, interests)
		return
	}

	httputil.NewResponse(
		statusCode,
		interests,
	).Send(c)
}
