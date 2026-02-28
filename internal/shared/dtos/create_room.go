package dtos

import "chat-app/internal/models"

type RoomResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Topic       string  `json:"topic"`
	IsPrivate   bool    `json:"is_private"`
	MaxMembers  int     `json:"max_members"`
	InviteCode  *string `json:"invite_code,omitempty"`
	Image       string  `json:"image"`
	ShareLink   string  `json:"share_link"`
	CreatedAt   string  `json:"created_at"`
	MemberCount int `json:"member_count"`
}

func MapToRoomResponse(room *models.Room, shareLink string) *RoomResponse {
	return &RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Topic:       room.Topic,
		IsPrivate:   room.IsPrivate,
		MaxMembers:  room.MaxMembers,
		InviteCode:  room.InviteCode,
		Image:       room.Image,
		ShareLink:   shareLink,
		CreatedAt:   room.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
