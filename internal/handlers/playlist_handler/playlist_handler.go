package playlisthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/internal/middleware"
	"github.com/ishantSikdar/mindo-server/pkg/utils/constant"
	"github.com/ishantSikdar/mindo-server/pkg/utils/http"
	"github.com/ishantSikdar/mindo-server/pkg/utils/route"
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
