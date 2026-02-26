package request

import "mime/multipart"

type CreateRoomRequest struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description" binding:"required"`
	Topic       string `form:"topic" binding:"required"`
	IsPrivate   bool   `form:"is_private"`
	MaxMembers  int    `form:"max_members" binding:"required"`
	Image       *multipart.FileHeader `form:"image" binding:"required"`
}