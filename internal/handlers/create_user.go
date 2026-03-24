package handlers

import (
	"chat-app/internal/services"
	"chat-app/internal/shared/config"
	"chat-app/internal/shared/logger"
	"chat-app/internal/shared/request"
	"chat-app/internal/shared/response"
	"chat-app/internal/shared/session"
	"chat-app/internal/shared/utils"
	"context"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateUser struct {
	Service services.CreateService
	Log     *logger.Logger
	Redis   *session.Store
	*utils.ImageUpload
}

func (s *CreateUser) CreateUserSession(c *gin.Context) {
	var user struct {
		Name string
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		s.Log.Error("invalid  data %v", err)
		response.BadRequest(c, nil, "Invalid  data")
		return
	}

	existingUser, err := s.Service.CheckUser(user.Name)

	if err != nil {

		if err.Error() == response.UserNotFound {
			response.BadRequest(c, nil, "user was not found")
			return
		}
		s.Log.Error("error%v", err)
		response.InternalServerError(c)
		return
	}

	data := &session.Data{
		ID:       existingUser.ID,
		UserName: existingUser.UserName,
	}

	newSessionID, err := s.Redis.CreateSession(c.Request.Context(), data)
	if err != nil {
		s.Log.Error("REDIS ERROR at line 103: %v | Context Err: %v", err, c.Request.Context().Err())
		response.InternalServerError(c)
		return
	}

	s.setSessionCookie(c ,newSessionID)
	response.Created(c, "created user session", data)

}

func (s *CreateUser) CreateUser(c *gin.Context) {

	var req request.CreateUser
	if err := c.ShouldBind(&req); err != nil {
		s.Log.Error("invalid form data %v", err)
		response.BadRequest(c, nil, "Invalid form data")
		return
	}

	resp, err := s.Service.CreateUser(req, "")
	
	if err != nil {
		s.Log.Error("creating user failed: %v", err)
		response.BadRequest(c, nil, "creating user failed")
		return
	}

	if req.ProfileImage!=nil{
		go s.handleAsyncImageUpload(resp.ID, req.ProfileImage)
	}

	data := &session.Data{
		ID:       resp.ID,
		UserName: resp.UserName,
	}

	newSessionID, err := s.Redis.CreateSession(c.Request.Context(), data)
	if err != nil {
		s.Log.Error("failed to create session: %v", err)
		response.InternalServerError(c)
		return
	}
	s.setSessionCookie(c ,newSessionID)

	response.Created(c, "created user successfully", data)

}

func (s *CreateUser) handleAsyncImageUpload(userID int, file *multipart.FileHeader) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    imgData, err := s.ImageUpload.UploadToCloudinary(ctx, file)
    if err != nil {
        s.Log.Error("Background Upload Failed for User %d: %v", userID, err)
        return
    }
    if err := s.Service.UpdateUserImage(userID, imgData); err != nil {
        s.Log.Error("Failed to update user %d with image URL: %v", userID, err)
    }
}

func (s *CreateUser) GetMe(c *gin.Context) {
	s.Log.Info("GetMe called")
	userID, exists := c.Get("user_id")
	if !exists {
		s.Log.Error("GetMe: user_id not found in context")
		response.Unauthorized(c, "unauthorized")
		return
	}

	user, err := s.Service.GetMe(userID.(int))
	if err != nil {
		s.Log.Error("GetMe: failed to get user: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, response.SuccessMsgFetched, user)
}

func NewUserHandler(srv services.CreateService, log *logger.Logger, redis *session.Store, cfg config.CloudinaryConfig) *CreateUser {
	imgUploader, err := utils.NewImageUpload(cfg, log)
	if err != nil {
		log.Error("Failed to initialize Cloudinary: %v", err)
		os.Exit(1)

	}
	return &CreateUser{
		Service:     srv,
		Log:         log,
		Redis:       redis,
		ImageUpload: imgUploader,
	}
}


func (s *CreateUser) setSessionCookie(c *gin.Context, sessionID string) {
    c.SetSameSite(http.SameSiteNoneMode)
    // secure := c.Request.TLS != nil
    c.SetCookie("session_id", sessionID, int(24*time.Hour.Seconds()), "/", "", false, true)
}