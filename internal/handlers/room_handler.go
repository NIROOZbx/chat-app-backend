package handlers

import (
	"chat-app/internal/services"
	"chat-app/internal/shared/config"
	"chat-app/internal/shared/logger"
	"chat-app/internal/shared/request"
	"chat-app/internal/shared/response"
	"chat-app/internal/shared/utils"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	Service services.Roomservice
	Log     *logger.Logger
	Cld     *utils.ImageUpload
}

func NewRoomHandler(srv services.Roomservice, log *logger.Logger, cfg config.CloudinaryConfig) *RoomHandler {
	imgUpload, err := utils.NewImageUpload(cfg, log)

	if err != nil {
		return nil
	}
	return &RoomHandler{
		Service: srv,
		Log:     log,
		Cld:     imgUpload,
	}
}

func (r *RoomHandler) CreateRoom(c *gin.Context) {
	var req request.CreateRoomRequest
	if err := c.ShouldBind(&req); err != nil {
		r.Log.Error("failed to bind request: %v", err)
		response.BadRequest(c, nil, "invalid form data")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.InternalServerError(c)
		return
	}

	imgURL, err := r.Cld.UploadToCloudinary(c.Request.Context(), req.Image)

	if err != nil {
		response.BadRequest(c, nil, "failed to upload image")
		return
	}

	data, err := r.Service.CreateRoom(c.Request.Context(), req, userID.(int), imgURL)

	if err != nil {
		r.Log.Error("error occurred %v", err)
		if err.Error() == "maximum 50 members allowed" {
			response.BadRequest(c, nil, err.Error())
			return
		}
		response.InternalServerError(c)
		return
	}

	response.Created(c, "Room "+response.SuccessMsgCreated, data)
}

func (r *RoomHandler) GetAllRooms(c *gin.Context) {

	rooms, err := r.Service.GetAllRooms()

	if err != nil {
		r.Log.Error("GetAllRooms: Failed to fetch: %v", err)
		response.InternalServerError(c)
	}

	response.OK(c, response.SuccessMsgFetched, rooms)

}

func (r *RoomHandler) GetJoinedRooms(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.InternalServerError(c)
		return
	}

	rooms, err := r.Service.GetJoinedRooms(userID.(int))
	if err != nil {
		r.Log.Error("GetJoinedRooms: Failed to fetch: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, response.SuccessMsgFetched, rooms)
}

func (r *RoomHandler) GetSingleRoom(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, nil, "Invalid room ID format")
		return
	}

	room, err := r.Service.GetSingleRoom(id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Log.Error("room not found %v", err.Error())
			response.NotFound(c, "Room not found")
			return
		}
		r.Log.Error("GetSingleRoom: Failed to fetch ID %d: %v", id, err)
		response.InternalServerError(c)
		return
	}

	onlineCount := r.Service.GetOnlineCount(id)
	response.OK(c, response.SuccessMsgFetched, gin.H{
		"room":         room,
		"online_count": onlineCount,
	})

}

func (r *RoomHandler) DeleteRoom(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.InternalServerError(c)
		return
	}

	if err := r.Service.DeleteRoom(id, userID.(int)); err != nil {
		if err == sql.ErrNoRows {
			r.Log.Error("room not found %v", err.Error())
			response.NotFound(c, "Room not found")
			return
		}

		if err.Error() == "unauthorized: only admins can delete rooms" {
			response.Forbidden(c, nil, err.Error())
			return
		}

		r.Log.Error("DeleteRoom: Failed to delete ID %d: %v", id, err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, response.SuccessMsgDeleted, nil)
}
