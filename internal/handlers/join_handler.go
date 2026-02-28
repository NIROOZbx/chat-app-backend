package handlers

import (
	"chat-app/internal/services"
	"chat-app/internal/shared/config"
	"chat-app/internal/shared/logger"
	"chat-app/internal/shared/response"
	"chat-app/internal/shared/session"
	"chat-app/internal/shared/utils"
	"github.com/gin-gonic/gin"
)

type JoinRoomHandler struct {
	Service services.JoinRoomService
	Log     *logger.Logger
}

func (s *JoinRoomHandler) JoinRoom(c *gin.Context) {

	userID := c.GetInt("user_id")
	roomID, ok := utils.ParseIDParam(c,"id")
	if !ok {
		return
	}

	resp, err := s.Service.JoinRoom(c.Request.Context(),userID, roomID)
	if err != nil {
		s.Log.Error("Join Room error: %v", err)
		if err.Error() == "user is already a member of this room" {
			response.BadRequest(c, nil, err.Error())
			return
		}
		response.InternalServerError(c)
		return
	}
	response.OK(c, "joined room successfully", resp)

}

func (s *JoinRoomHandler) LeaveRoom(c *gin.Context) {

	userID := c.GetInt("user_id")
	userName:=c.GetString("user_name")

	roomID, ok := utils.ParseIDParam(c, "id")

	if !ok {
		s.Log.Error("error in converting ")
		return
	}
	err := s.Service.LeaveRoom(c.Request.Context(),userName,roomID, userID)
	if err != nil {
		s.Log.Error("error %v", err)
		if err.Error() == response.UserNotMember {
			response.NotFound(c, "You are not a member of this room")
			return
		}
		s.Log.Error("LeaveRoom: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, "room left successfully", nil)

}

func (s *JoinRoomHandler)JoinPrivateGroup(c *gin.Context){

	userID := c.GetInt("user_id")
	

	id:=c.Param("inviteCode")

	room,err:=s.Service.JoinPrivateRoom(c.Request.Context(),userID,id)
	if err!=nil{
		s.Log.Error("join error %v",err)
		response.BadRequest(c,nil,"failed to join")
		return
	}

	response.OK(c,"room joined successfully",room)


}

func NewJoinRoomHandler(srv services.JoinRoomService, log *logger.Logger, cfg config.CloudinaryConfig, store *session.Store) *JoinRoomHandler {


	return &JoinRoomHandler{
		Service:     srv,
		Log:         log,

	}
}
