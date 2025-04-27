package handlers

import (
	authhandler "github.com/easc01/mindo-server/internal/handlers/auth_handler"
	interesthandler "github.com/easc01/mindo-server/internal/handlers/interest_handler"
	userhandler "github.com/easc01/mindo-server/internal/handlers/user_handler"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitREST() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
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
		userhandler.RegisterAppUserRoutes(apiRg)
		userhandler.RegisterAdminUserRoutes(apiRg)
		interesthandler.RegisterInterest(apiRg)
	}
}
