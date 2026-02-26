package handlers

import (
	"chat-app/internal/services"
	"chat-app/internal/shared/config"
	"chat-app/internal/shared/logger"
	"chat-app/internal/shared/request"
	"chat-app/internal/shared/response"
	"chat-app/internal/shared/session"
	"chat-app/internal/shared/utils"
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
	}

	data := &session.Data{
		UserID:   existingUser.ID,
		UserName: existingUser.UserName,
	}

	newSessionID, err := s.Redis.CreateSession(c.Request.Context(), data)
	if err != nil {
		s.Log.Error("failed to create session: %v", err)
		response.InternalServerError(c)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("session_id", newSessionID, int(24*time.Hour.Seconds()), "/", "", true, true)

	response.Created(c, "created user successfully", data)

}

func (s *CreateUser) CreateUser(c *gin.Context) {

	var req request.CreateUser
	if err := c.ShouldBind(&req); err != nil {
		s.Log.Error("invalid form data %v", err)
		response.BadRequest(c, nil, "Invalid form data")
		return
	}

	var imageURL string
	if req.ProfileImage != nil {
		imgData, err := s.ImageUpload.UploadToCloudinary(c.Request.Context(), req.ProfileImage)
		if err != nil {
			s.Log.Error("upload error %v", err)
			response.InternalServerError(c)
			return
		}
		imageURL = imgData
	}

	resp, err := s.Service.CreateUser(req, imageURL)

	if err != nil {
		s.Log.Error("creating user failed: %v", err)
		response.BadRequest(c, nil, "creating user failed")
		return
	}

	data := &session.Data{
		UserID:   resp.ID,
		UserName: resp.UserName,
	}

	newSessionID, err := s.Redis.CreateSession(c.Request.Context(), data)
	if err != nil {
		s.Log.Error("failed to create session: %v", err)
		response.InternalServerError(c)
		return
	}
	c.SetSameSite(http.SameSiteNoneMode)

	c.SetCookie("session_id", newSessionID, int(24*time.Hour.Seconds()), "/", "", true, true)

	response.Created(c, "created user successfully", data)

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
