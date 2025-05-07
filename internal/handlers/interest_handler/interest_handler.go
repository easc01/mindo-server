package interesthandler

import (
	"net/http"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	interestservice "github.com/easc01/mindo-server/internal/services/interest_service"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/message"
	networkutil "github.com/easc01/mindo-server/pkg/utils/network_util"
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
	req, ok := networkutil.GetRequestBody[dto.UpsertInterestDTO](c)
	if !ok {
		return
	}

	user, ok := middleware.GetUser(c)
	if user.AdminUser == nil || !ok {
		logger.Log.Errorf(message.NullAppUserContext)
		networkutil.NewErrorResponse(
			http.StatusInternalServerError,
			message.NullAdminUserContext,
			nil,
		).Send(c)
		return
	}

	statusCode, upsertErr := interestservice.UpsertIntoMasterInterest(
		c,
		req.Interests,
		user.AdminUser.UserID.String(),
	)

	if upsertErr != nil {
		networkutil.NewErrorResponse(
			statusCode,
			upsertErr.Error(),
			nil,
		).Send(c)
		return
	}

	networkutil.NewResponse(
		statusCode,
		"interest master list updated",
	).Send(c)
}

func getMasterInterestListHandler(c *gin.Context) {
	interests, statusCode, intErr := interestservice.GetMasterInterestList(c)

	if intErr != nil {
		networkutil.NewErrorResponse(statusCode, message.SomethingWentWrong, interests)
		return
	}

	networkutil.NewResponse(
		statusCode,
		interests,
	).Send(c)
}
