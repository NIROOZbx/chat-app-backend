package utils

import (
	"chat-app/internal/shared/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseIDParam(c *gin.Context, paramName string) (int, bool) {
	val := c.Param(paramName)
	id, err := strconv.Atoi(val)
	if err != nil {
		response.BadRequest(c, nil, "Invalid "+paramName+" format")
		return 0, false
	}
	return id, true
}