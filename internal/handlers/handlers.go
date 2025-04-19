package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/internal/constants"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
)

func InitREST() {
	r := gin.Default()
	registerRoutes(&r.RouterGroup)

	routerErr := r.Run(":8080")

	if routerErr != nil {
		logger.Log.Error("failed to start router")
	}
}

func registerRoutes(rg *gin.RouterGroup) {
	apiRg := rg.Group(constants.Api)

	{
		RegisterAuth(apiRg)
		RegisterPlaylist(apiRg)
		RegisterUserRoutes(apiRg)
	}

	logger.Log.Info("registered routes")
}
