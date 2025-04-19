package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/internal/constants"
	"github.com/ishantSikdar/mindo-server/internal/middleware"
	"github.com/ishantSikdar/mindo-server/pkg/utils"
)

func RegisterPlaylist(rg *gin.RouterGroup) {
	playlistRg := rg.Group(constants.Playlists, middleware.AuthMiddleware())
	{
		playlistRg.POST(utils.Blank, createPlaylistHandler)
		playlistRg.GET(utils.Blank, getAllPlaylistsHandler)
		playlistRg.GET(utils.IdParam, getPlaylistByIdHandler)
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
