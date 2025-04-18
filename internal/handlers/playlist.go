package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ishantSikdar/mindo-server/internal/models"
	"github.com/ishantSikdar/mindo-server/pkg/db"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
)

const PLAYLIST = "/playlists"

func RegisterPlaylist(rg *gin.RouterGroup) {
	playlistRg := rg.Group("/playlists")
	{
		playlistRg.POST("", createPlaylistHandler)
		playlistRg.GET("", getAllPlaylistsHandler)
		playlistRg.GET("/:id", getPlaylistByIdHandler)
	}
	logger.Log.Info("Registered playlist routes")
}

func createPlaylistHandler(c *gin.Context) {
	params := models.CreatePlaylistParams{
		Name:         sql.NullString{String: "My Static Playlist", Valid: true},
		Description:  sql.NullString{String: "This is a statically created playlist", Valid: true},
		ThumbnailUrl: sql.NullString{String: "http://example.com/static_thumbnail.jpg", Valid: true},
		UpdatedBy:    uuid.NullUUID{UUID: uuid.MustParse("a93a5c8d-dc57-4c1b-8a4e-432f3bba1b32"), Valid: true}, // Static UUID
	}

	playlist, err := db.Queries.CreatePlaylist(context.Background(), params)
	if err != nil {
		logger.Log.Error("Error saving playlist", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong when saving Playlist",
		})
	}

	c.JSON(http.StatusOK, playlist)
}

func getAllPlaylistsHandler(c *gin.Context) {
	playlists, err := db.Queries.GetAllPlaylists(context.Background())

	if err != nil {
		logger.Log.Error("Failed to get playlists", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong while fetching playlists",
		})
	}

	c.JSON(http.StatusOK, playlists)
}

func getPlaylistByIdHandler(c *gin.Context) {
	id := c.Param("id")

	parsedId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid playlist ID format",
		})
		return
	}

	playlist, err := db.Queries.GetPlaylistByID(context.Background(), parsedId)

	if err != nil {
		logger.Log.Error("Failed to get playlist", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong while fetching playlist",
		})
	}

	c.JSON(http.StatusOK, playlist)
}
