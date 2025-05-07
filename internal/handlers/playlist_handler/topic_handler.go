package playlisthandler

import (
	"net/http"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	playlistservice "github.com/easc01/mindo-server/internal/services/playlist_service"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	networkutil "github.com/easc01/mindo-server/pkg/utils/network_util"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterTopic(rg *gin.RouterGroup) {
	topicRg := rg.Group(route.Topics)

	{
		topicRg.GET(
			constant.IdParam+"/videos",
			middleware.RequireRole(models.UserTypeAppUser, models.UserTypeAdminUser),
			getTopicVideosHandler,
		)
	}
}

func getTopicVideosHandler(c *gin.Context) {
	topicId := c.Param("id")
	videoId := c.Query("videoId")

	parsedTopicId, err := uuid.Parse(topicId)
	if err != nil {
		networkutil.NewErrorResponse(
			http.StatusBadRequest,
			"invalid topic id",
			err.Error(),
		).Send(c)
		return
	}

	videos, statusCode, err := playlistservice.GetVideosByTopicId(c, parsedTopicId, videoId)
	if err != nil {
		networkutil.NewErrorResponse(
			statusCode,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	networkutil.NewResponse(
		statusCode,
		videos,
	).Send(c)
}
