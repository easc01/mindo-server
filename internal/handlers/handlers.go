package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	apiRg := rg.Group("/api")
	RegisterPlaylist(apiRg)
	logger.Log.Info("Registered routes")
}
