package httputil

import (
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response
type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
}

// NewResponse creates a new response object
func NewResponse(status int, data interface{}) *Response {
	return &Response{
		StatusCode: status,
		Message:    constant.Blank,
		Data:       data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(status int, errMessage string, data any) *Response {
	return &Response{
		StatusCode: status,
		Message:    errMessage,
		Data:       data,
	}
}

// Send writes the response to the client
func (r *Response) Send(c *gin.Context) {
	c.JSON(r.StatusCode, r)
}
