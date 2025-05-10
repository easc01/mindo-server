package communityhandler

import (
	"net/http"
	"time"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	communityservice "github.com/easc01/mindo-server/internal/services/community_service"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	networkutil "github.com/easc01/mindo-server/pkg/utils/network_util"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterMessages(rg *gin.RouterGroup) {
	messageRg := rg.Group(route.Messages)

	{
		messageRg.GET(
			constant.Blank,
			middleware.RequireRole(models.UserTypeAppUser),
			messageHistoryPageHandler,
		)
	}
}

func messageHistoryPageHandler(c *gin.Context) {
	communityID := c.Query("communityId")
	lastMsgTime := c.Query("lastMessageTime")

	parsedCommID, err := uuid.Parse(communityID)

	if err != nil {
		networkutil.NewErrorResponse(
			http.StatusBadRequest,
			"invalid community id",
			nil,
		).Send(c)
		return
	}

	parsedTime, err := time.Parse(constant.TimeLayout, lastMsgTime)
	if err != nil {
		parsedTime = time.Now()
	}

	userMessages, statusCode, err := communityservice.GetMessageHistoryPage(
		c,
		parsedCommID,
		parsedTime,
	)

	if err != nil {
		networkutil.NewErrorResponse(
			statusCode,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	networkutil.NewResponse(
		http.StatusAccepted,
		userMessages,
	).Send(c)
}
