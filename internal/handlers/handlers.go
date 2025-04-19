package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/internal/constants"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
)

func InitREST() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	registerRoutes(&r.RouterGroup)

	routerErr := r.Run(":8080")

	if routerErr != nil {
		logger.Log.Errorf("failed to start router, %s", routerErr)
	}
}

func registerRoutes(rg *gin.RouterGroup) {
	apiRg := rg.Group(constants.Api)

	{
		RegisterAuth(apiRg)
		RegisterPlaylist(apiRg)
		RegisterUserRoutes(apiRg)
	}
}
