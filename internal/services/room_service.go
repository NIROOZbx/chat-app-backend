package services

import (
	"chat-app/internal/models"
	"chat-app/internal/repositories"
	"chat-app/internal/shared/hub"
	"chat-app/internal/shared/request"
	"context"
	"errors"
)

type Roomservice interface {
	CreateRoom(ctx context.Context, req request.CreateRoomRequest, creatorID int, img string) (*models.Room, error)
	GetAllRooms() ([]models.Room, error)
	GetJoinedRooms(userID int) ([]models.Room, error)
	GetSingleRoom(id int) (*models.Room, error)
	DeleteRoom(id int) error
	GetOnlineCount(roomID int) int
}

type roomService struct {
	Repo repositories.RoomRepository
	Hub  *hub.Manager
}

func (s *roomService) CreateRoom(ctx context.Context, req request.CreateRoomRequest, creatorID int, img string) (*models.Room, error) {

	if req.MaxMembers > 50 {
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
	err := s.Repo.CreateRoom(room, creatorID)
	if err != nil {
		return nil, err
	}

	return room, nil
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

func (r *roomService) DeleteRoom(id int) error {
	return r.Repo.DeleteRoom(id)
}

func (s *roomService) GetOnlineCount(roomID int) int {
	return s.Hub.GetOnlineCount(roomID)
}
func NewRoomService(repo repositories.RoomRepository, hub *hub.Manager) Roomservice {
	return &roomService{Repo: repo, Hub: hub}
}
