package logger

import "github.com/gin-gonic/gin"

// GinError sets error details in the response headers and terminates the request with the associated status code.
func (e *Error) GinError(c *gin.Context) {
	c.Header("x-error", e.Message)
	c.Header("x-error-id", e.UUID)
	c.AbortWithStatus(e.Status)
}
