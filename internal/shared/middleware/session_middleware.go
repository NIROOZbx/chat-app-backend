package middleware

import (
	"chat-app/internal/shared/response"
	"chat-app/internal/shared/session"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SessionMiddleware(store *session.Store) gin.HandlerFunc {

	return func(c *gin.Context) {

		sessionID, err := c.Cookie("session_id")

		if err != nil || sessionID == "" {
			response.NotFound(c, "cookie not found")
			fmt.Println(sessionID)
			c.Abort()
			return
		}

		data, err := store.Get(c.Request.Context(), sessionID)
		if err != nil {
			response.InternalServerError(c)
			c.Abort()
			return
		}
		if data == nil {
			c.SetCookie("session_id", "", -1, "/", "", true, true) // clear stale cookie
			response.BadRequest(c, nil, "Session expired — join again")
			c.Abort()
			return
		}
		c.Set("user_id", data.UserID)
		c.Set("user_name", data.UserName)
		c.Next()

	}

}
