package httputil

import (
	"net/http"

	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/gin-gonic/gin"
)

func GetRequestBody[T any](c *gin.Context) (T, bool) {
	var reqBody T
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Log.Error(message.InvalidRequestBody)

		NewErrorResponse(
			http.StatusBadRequest,
			message.InvalidRequestBody,
			err.Error(),
		).Send(c)
		c.Abort()
		return reqBody, false
	}
	return reqBody, true
}
