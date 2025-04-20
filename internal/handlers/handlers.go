package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	authhandler "github.com/ishantSikdar/mindo-server/internal/handlers/auth_handler"
	playlisthandler "github.com/ishantSikdar/mindo-server/internal/handlers/playlist_handler"
	userhandler "github.com/ishantSikdar/mindo-server/internal/handlers/user_handler"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/utils/route"
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
	apiRg := rg.Group(route.Api)

	{
		authhandler.RegisterAuth(apiRg)
		playlisthandler.RegisterPlaylist(apiRg)
		userhandler.RegisterUserRoutes(apiRg)
	}
}
