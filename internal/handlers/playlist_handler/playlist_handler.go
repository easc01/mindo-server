package playlisthandler

import (
	"net/http"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	playlistservice "github.com/easc01/mindo-server/internal/services/playlist_service"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	httputil "github.com/easc01/mindo-server/pkg/utils/http_util"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterPlaylists(rg *gin.RouterGroup) {
	playlistRg := rg.Group(route.Playlists)

	{
		playlistRg.POST(
			constant.Blank,
			middleware.RequireRole(models.UserTypeAdminUser),
			createPlaylistHandler,
		)

		playlistRg.GET(
			constant.Blank,
			middleware.RequireRole(
				models.UserTypeAdminUser,
				models.UserTypeAppUser,
			), // TODO, app user to be used later to map Playlist by user interests
			getAllPlaylistPreviews,
		)

		playlistRg.GET(
			constant.IdParam,
			middleware.RequireRole(models.UserTypeAppUser, models.UserTypeAdminUser),
			getPlaylistByIdHandler,
		)
	}
}

func createPlaylistHandler(c *gin.Context) {
	req, ok := httputil.GetRequestBody[dto.CreatePlaylistRequest](c)
	if !ok {
		return
	}

	var validTopics []string

	// Filter valid topics
	for _, topic := range req.Topics {
		if topic != constant.Blank {
			validTopics = append(validTopics, topic)
		}
	}

	// Ensure valid topics exist
	if len(validTopics) == 0 {
		httputil.NewErrorResponse(
			http.StatusBadRequest,
			"no valid topics provided",
			nil,
		).Send(c)
		return
	}

	req.Topics = validTopics

	user, ok := middleware.GetUser(c)
	if user.AdminUser == nil || !ok {
		logger.Log.Errorf(message.NullAdminUserContext)
		httputil.NewErrorResponse(
			http.StatusInternalServerError,
			message.SomethingWentWrong,
			message.NullAdminUserContext,
		).Send(c)
		return
	}

	playlistDetails, statusCode, err := playlistservice.ProcessPlaylistCreationByAdmin(
		c,
		req,
		user.AdminUser.UserID,
	)

	if err != nil {
		logger.Log.Errorf("failed to process playlist creation by admin, %s", user.AdminUser.UserID)
		httputil.NewErrorResponse(
			statusCode,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	httputil.NewResponse(
		http.StatusCreated,
		playlistDetails,
	).Send(c)
}

func getAllPlaylistPreviews(c *gin.Context) {
	playlists, statusCode, err := playlistservice.GetAllPlaylistPreviews(c)
	if err != nil {
		logger.Log.Error("failed to get playlist previews")
		httputil.NewErrorResponse(
			statusCode,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	httputil.NewResponse(
		statusCode,
		playlists,
	).Send(c)
}

func getPlaylistByIdHandler(c *gin.Context) {
	playlistId := c.Param("id")

	parsedPlaylistId, err := uuid.Parse(playlistId)
	if err != nil {
		httputil.NewErrorResponse(
			http.StatusBadRequest,
			message.InvalidUserID,
			err.Error(),
		).Send(c)
		return
	}

	playlistData, statusCode, err := playlistservice.GetPlaylistWithTopics(c, parsedPlaylistId)
	if err != nil {
		httputil.NewErrorResponse(
			statusCode,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	httputil.NewResponse(
		statusCode,
		playlistData,
	).Send(c)
}
