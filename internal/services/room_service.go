package services

import (
	"chat-app/internal/models"
	"chat-app/internal/repositories"
	"chat-app/internal/shared/dtos"
	"chat-app/internal/shared/hub"
	"chat-app/internal/shared/request"
	"chat-app/internal/shared/utils"
	"context"
	"errors"
	"fmt"
	"os"
)

type Roomservice interface {
	CreateRoom(ctx context.Context, req request.CreateRoomRequest, creatorID int, img string) (*dtos.RoomResponse, error)
	GetAllRooms() ([]models.Room, error)
	GetJoinedRooms(userID int) ([]models.Room, error)
	GetSingleRoom(id int) (*models.Room, error)
	DeleteRoom(roomID, userID int) error
	GetOnlineCount(roomID int) int
}

type roomService struct {
	Repo repositories.RoomRepository
	Hub  *hub.Manager
}

func (s *roomService) CreateRoom(ctx context.Context, req request.CreateRoomRequest, creatorID int, img string) (*dtos.RoomResponse, error) {

	if req.MaxMembers <= 0 || req.MaxMembers > 50 {
		return nil, errors.New("maximum 50 members allowed")
	}

	room := &models.Room{
		Name:        req.Name,
		Description: req.Description,
		Topic:       req.Topic,
		IsPrivate:   req.IsPrivate,
		MaxMembers:  req.MaxMembers,
		Image:       img,
	}

	var link string

	if req.IsPrivate {
		room.InviteCode = utils.GenerateInviteCode(8)
		if room.InviteCode == nil {
			return nil, errors.New("failed to generate invite code")
		}
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			return nil, errors.New("FRONTEND_URL environment variable not set")
		}
		link = fmt.Sprintf("%s/join/%s", frontendURL, *room.InviteCode)
	}
	err := s.Repo.CreateRoom(room, creatorID)
	if err != nil {
		return nil, err
	}

	return dtos.MapToRoomResponse(room, link), nil
}

func (r *roomService) GetAllRooms() ([]models.Room, error) {
	return r.Repo.GetAllRooms()
}

func (s *roomService) GetJoinedRooms(userID int) ([]models.Room, error) {
	return s.Repo.GetJoinedRooms(userID)
}

func (r *roomService) GetSingleRoom(id int) (*models.Room, error) {
	return r.Repo.GetRoomById(id)
}

func (r *roomService) DeleteRoom(roomID, userID int) error {
	role, err := r.Repo.GetUserRole(roomID, userID)
	if err != nil {
		return err
	}

	if role != "admin" {
		return errors.New("unauthorized: only admins can delete rooms")
	}

	return r.Repo.DeleteRoom(roomID)
}

func (s *roomService) GetOnlineCount(roomID int) int {
	return s.Hub.GetOnlineCount(roomID)
}
func NewRoomService(repo repositories.RoomRepository, hub *hub.Manager) Roomservice {
	return &roomService{Repo: repo, Hub: hub}
}
