package httputil

import (
	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response
type Response struct {
	Status int `json:"-"`
	Data   any `json:"data"`
}

// ErrorResponse contains detailed error information
type ErrorResponse struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// NewResponse creates a new response object
func NewResponse(status int, data interface{}) *Response {
	return &Response{
		Status: status,
		Data:   data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(status int, errMessage string, details any) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: errMessage,
		Details: details,
	}
}

// Send writes the response to the client
func (r *Response) Send(c *gin.Context) {
	c.JSON(r.Status, r)
}

// Send writes the response to the client
func (r *ErrorResponse) Send(c *gin.Context) {
	c.JSON(r.Status, r)
}
