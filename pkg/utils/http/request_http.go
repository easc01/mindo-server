package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/utils/message"
)

func GetRequestBody[T any](c *gin.Context) (T, bool) {
	var reqBody T
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Log.Error(message.InvalidRequestBody)
		NewErrorResponse(
			StatusBadRequest,
			message.InvalidRequestBody,
			nil,
		).Send(c)
		c.Abort()
		return reqBody, false
	}
	return reqBody, true
}
