package dtos

import "time"

type MessageResponse struct {
    ID        int       `json:"id"`
    RoomID    int       `json:"room_id"`
    UserID    int       `json:"user_id"`
    UserName  string    `json:"user_name"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}