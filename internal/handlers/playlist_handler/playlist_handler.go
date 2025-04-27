package playlisthandler

import (
	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
)

func RegisterPlaylists(rg *gin.RouterGroup) {
	playlistRg := rg.Group(route.Playlists)

	{
		playlistRg.POST(
			constant.Blank,
			middleware.RequireRole(models.UserTypeAdminUser),
			CreatePlaylistHandler,
		)

		playlistRg.GET(
			constant.Blank,
			middleware.RequireRole(models.UserTypeAppUser),
			GetAppUserInterestedPlaylistsHandler,
		)

		playlistRg.GET(
			constant.IdParam,
			middleware.RequireRole(models.UserTypeAppUser),
			GetPlaylistByIdHandler,
		)

		playlistRg.GET(
			constant.IdParam + "/videos",
			middleware.RequireRole(models.UserTypeAppUser),
			GetPlaylistTopicVideosHandler,
		)
	}
}

func CreatePlaylistHandler(c *gin.Context) {

}

func GetAppUserInterestedPlaylistsHandler(c *gin.Context) {

}

func GetPlaylistByIdHandler(c *gin.Context) {

}

func GetPlaylistTopicVideosHandler(c *gin.Context) {

}
