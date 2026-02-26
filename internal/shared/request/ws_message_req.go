package request

type WSMessage struct {
    Type    string `json:"type"`
    Content string `json:"content"`
}