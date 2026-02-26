package request

import "mime/multipart"

type CreateUser struct {
	UserName     string `form:"user_name" binding:"required"`
	ProfileImage *multipart.FileHeader `form:"profile_image"`
}