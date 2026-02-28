package response

import (
	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(StatusCreated, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}
func BadRequest(c *gin.Context, data interface{}, message string) {

	c.JSON(StatusBadRequest, gin.H{
		"success": false,
		"data":    data,
		"error":   message,
	})
}

func InternalServerError(c *gin.Context) {
	c.JSON(StatusInternalServerError, gin.H{
		"success": false,
		"data":    nil,
		"error":   "An unexpected server error occurred. Please try again later.",
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(StatusUnauthorized, gin.H{
		"success": false,
		"data":    nil,
		"error":   "An unexpected server error occurred. Please try again later.",
	})

}

func NotFound(c *gin.Context, message string) {
	c.JSON(StatusNotFound, gin.H{
		"success": false,
		"error":   message,
	})
}

func Conflict(c *gin.Context, message string) {
	c.JSON(StatusConflict, gin.H{
		"success": false,
		"error":   message,
	})
}

func Forbidden(c *gin.Context, data interface{}, message string) {
	c.JSON(StatusForbidden, gin.H{
		"success": false,
		"data":    data,
		"error":   message,
	})
}
