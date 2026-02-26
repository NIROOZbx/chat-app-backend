package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {

	c := cors.New(cors.Config{
		AllowOrigins:     []string{"https://chat-app-frontend-pink-phi.vercel.app", "http://localhost:5174"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie"},
	})

	return c

}
