package handlers

import (
	"chat-app/internal/services"
	"chat-app/internal/shared/logger"
	"chat-app/internal/shared/request"
	"chat-app/internal/shared/response"
	"chat-app/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
    Service services.MessageService
    Log     *logger.Logger
}

func NewMessageHandler(srv services.MessageService, log *logger.Logger) *MessageHandler {
    return &MessageHandler{Service: srv, Log: log}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
    userID   := c.GetInt("user_id")
    userName := c.GetString("user_name")

    roomID, ok := utils.ParseIDParam(c, "id")
    if !ok {
        return
    }

    var req request.SendMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.Log.Error("invalid body: %v", err)
        response.BadRequest(c, nil, "content is required")
        return
    }

    msg, err := h.Service.SendMessage(c.Request.Context(), roomID, userID, userName, req.Content)
    if err != nil {
        h.Log.Error("SendMessage: %v", err)
        response.InternalServerError(c)
        return
    }

    response.Created(c, "message sent", msg)
}


func (h *MessageHandler) GetMessages(c *gin.Context) {
    roomID, ok := utils.ParseIDParam(c, "id")
    if !ok {
        return
    }

	var req request.PaginationInput
    if err:=c.ShouldBindQuery(&req);err!=nil{
		h.Log.Error("error in binding%v",err)
		response.BadRequest(c,nil,"error in binding")
		return
	}


    messages, err := h.Service.GetMessages(c.Request.Context(),roomID, req.Limit, req.Page)
    if err != nil {
        h.Log.Error("GetMessages: %v", err)
        response.InternalServerError(c)
        return
    }

    response.OK(c, "messages fetched", messages)
}