package playlisthandler

import (
	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/http"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
)

func RegisterPlaylist(rg *gin.RouterGroup) {
	playlistRg := rg.Group(route.Playlists, middleware.AuthMiddleware())
	{
		playlistRg.POST(constant.Blank, createPlaylistHandler)
		playlistRg.GET(constant.Blank, getAllPlaylistsHandler)
		playlistRg.GET(constant.IdParam, getPlaylistByIdHandler)
	}
}

func createPlaylistHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "comming soon",
	})
}

func getAllPlaylistsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "comming soon",
	})
}

func getPlaylistByIdHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "comming soon",
	})
}
