package communityhandler

import (
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
)

func RegisterCommunity(rg *gin.RouterGroup) {
	communityRg := rg.Group(route.Admin)

	{
		communityRg.POST(constant.Blank, createCommunity)
		communityRg.POST("/join/:communityId", joinCommunity)
	}
}

func createCommunity(c *gin.Context) {
	
}

func joinCommunity(c *gin.Context) {

}
