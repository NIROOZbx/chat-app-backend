package dtos

import "chat-app/internal/models"

type JoinRoomResponse struct {
	UserID          int    `json:"user_id"`
	UserName        string `json:"user_name"`
	ProfileImage    string `json:"profile_image"`
	RoomID          int    `json:"room_id"`
	RoomName        string `json:"room_name"`
	RoomTopic       string `json:"room_topic"`
	RoomDescription string `json:"room_description"`
}

func MapToJoinResponse(user *models.User, room *models.Room) *JoinRoomResponse {
	return &JoinRoomResponse{
		UserID:          user.ID,
		UserName:        user.UserName,
		ProfileImage:    user.ProfileImage,
		RoomID:          room.ID,
		RoomName:        room.Name,
		RoomTopic:       room.Topic,
		RoomDescription: room.Description,
	}
}
