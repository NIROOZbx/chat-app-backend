package handlers

import (
	"chat-app/internal/services"
	"chat-app/internal/shared/logger"
	"chat-app/internal/shared/request"
	"chat-app/internal/shared/response"
	"chat-app/internal/shared/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSHandler struct {
	service        services.WSService
	log            *logger.Logger
	messageService services.MessageService
}

func (ws *WSHandler) CreateWSConn(c *gin.Context) {
	ws.log.Error("CreateWSConn called")

	userID := c.GetInt("user_id")
	userName := c.GetString("user_name")

	if userID == 0 || userName == "" {
		response.BadRequest(c, nil, "missing session")
		return
	}

	roomID, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.BadRequest(c, nil, "error in connecting to socket")
		return
	}
	defer conn.Close()

	ws.service.Connect(userName, roomID, userID, conn)
	defer ws.service.Disconnect(userName, roomID, userID)

	for {

		var msg request.WSMessage

		if err := conn.ReadJSON(&msg); err != nil {
			ws.log.Error("user %d disconnected: %v", userID, err)
			break
		}

		switch msg.Type {
		case "message.send":
			_, err := ws.messageService.SendMessage(context.Background(),
				roomID,
				userID,
				userName,
				msg.Content)

			if err != nil {
				ws.log.Error("SendMessage: %v", err)
			}
		case "room.typing":
			ws.service.UserTyping(context.Background(), userName, userID, roomID)

		}

	}
}

func NewWSConn(service services.WSService, log *logger.Logger, messageService services.MessageService) *WSHandler {
	return &WSHandler{
		messageService: messageService,
		service:        service,
		log:            log,
	}
}
