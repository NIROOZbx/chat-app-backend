package request

type SendMessageRequest struct {
    Content string `json:"content" binding:"required"`
}