package routes

import (
	"chat-app/internal/handlers"
	"chat-app/internal/shared/middleware"
	"chat-app/internal/shared/session"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, roomHandle *handlers.RoomHandler, joinHandle *handlers.JoinRoomHandler, userHandle *handlers.CreateUser, store *session.Store, ws *handlers.WSHandler, messageHandler *handlers.MessageHandler) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Chat Backend is running",
		})
	})
	v1 := r.Group("/api/v1")
	{

		v1.POST("/create-user", userHandle.CreateUser)
		v1.POST("/auth-user", userHandle.CreateUserSession)
		v1.GET("/me", middleware.SessionMiddleware(store), userHandle.GetMe)
		room := v1.Group("/rooms")
		{
			protected := room.Group("", middleware.SessionMiddleware(store))
			protected.POST("/join/:id", joinHandle.JoinRoom)
			protected.POST("/join/private/:inviteCode", joinHandle.JoinPrivateGroup)
			protected.DELETE("/leave/:id", joinHandle.LeaveRoom)
			protected.POST("/", roomHandle.CreateRoom)
			protected.GET("/joined", roomHandle.GetJoinedRooms)
			protected.GET("/ws/:id", ws.CreateWSConn)
			protected.POST("/:id/messages", messageHandler.SendMessage)
			protected.GET("/:id/messages", messageHandler.GetMessages)

			room.GET("/", roomHandle.GetAllRooms)
			room.GET("/:id", roomHandle.GetSingleRoom)
			protected.DELETE("/:id", roomHandle.DeleteRoom)
		}

	}

}
