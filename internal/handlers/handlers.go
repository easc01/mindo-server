package handlers

import (
	"net/http"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/easc01/mindo-server/internal/config"
	authhandler "github.com/easc01/mindo-server/internal/handlers/auth_handler"
	communityhandler "github.com/easc01/mindo-server/internal/handlers/community_handler"
	interesthandler "github.com/easc01/mindo-server/internal/handlers/interest_handler"
	playlisthandler "github.com/easc01/mindo-server/internal/handlers/playlist_handler"
	quizhandler "github.com/easc01/mindo-server/internal/handlers/quiz_handler"
	userhandler "github.com/easc01/mindo-server/internal/handlers/user_handler"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitREST() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://app.mindo.easc01.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	registerRoutes(&r.RouterGroup)
	registerWebSockets(r)

	r.GET("/", func(c *gin.Context) {
		hostInfo, _ := host.Info()
		cpuInfo, _ := cpu.Info()
		memInfo, _ := mem.VirtualMemory()
		diskInfo, _ := disk.Usage("/")

		c.JSON(http.StatusOK, gin.H{
			"host": hostInfo,
			"cpu":  cpuInfo,
			"memory":  memInfo,
			"disk": diskInfo,
		})
	})

	routerErr := r.Run(":" + config.GetConfig().AppPort)

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
		playlisthandler.RegisterPlaylists(apiRg)
		playlisthandler.RegisterTopic(apiRg)
		communityhandler.RegisterCommunity(apiRg)
		communityhandler.RegisterMessages(apiRg)
		quizhandler.RegisterQuiz(apiRg)
	}
}

func registerWebSockets(r *gin.Engine) {
	r.GET("/chat", func(ctx *gin.Context) {
		communityhandler.HandleRoomChatWS(ctx.Writer, ctx.Request)
	})
}
