package logger

import "github.com/gin-gonic/gin"

// GinError is a helper function to return an error to the client using Gin framework context.
// It sets the headers x-error and x-error-id with the error message and UUID respectively and sets the status code.
func (e *Error) GinError(c *gin.Context) {
	c.Header("x-error", e.Message)
	c.Header("x-error-id", e.UUID)
	c.AbortWithStatus(e.Status)
}
