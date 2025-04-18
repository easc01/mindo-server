package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
)

func RegisterRoutes(r *gin.Engine) {
	RegisterPlaylist(r)
	logger.Log.Info("Registered routes")
}
