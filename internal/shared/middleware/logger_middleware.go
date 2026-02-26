package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		log.Printf("[API] | %3d | %13v | %-7s | %s",
			c.Writer.Status(),
			latency,
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}
