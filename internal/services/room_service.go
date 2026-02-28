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

	"github.com/redis/go-redis/v9"
)

type Roomservice interface {
	CreateRoom(ctx context.Context, req request.CreateRoomRequest, creatorID int, img string) (*dtos.RoomResponse, error)
	GetAllRooms(ctx context.Context) ([]models.Room, error)
	GetJoinedRooms(ctx context.Context,userID int)([]models.Room, error)
	GetSingleRoom(id int) (*models.Room, error)
	DeleteRoom(roomID, userID int) error
	GetOnlineCount(roomID int) int
	GetUserRole(roomID, userID int) (string, error)
}

type roomService struct {
	Repo repositories.RoomRepository
	Hub  *hub.Manager
	Redis *redis.Client
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

	key := fmt.Sprintf("room:%d:member_count", room.ID)
    s.Redis.Set(ctx, key, 1, 0) 

	return dtos.MapToRoomResponse(room, link), nil
}

func (s *roomService) GetAllRooms(ctx context.Context) ([]models.Room, error) {
	rooms, err := s.Repo.GetAllRooms()
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %w", err)
	}

	if err := s.enrichWithMemberCount(ctx, rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *roomService) GetJoinedRooms(ctx context.Context, userID int) ([]models.Room, error) {
	rooms, err := s.Repo.GetJoinedRooms(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get joined rooms: %w", err)
	}

	if err := s.enrichWithMemberCount(ctx, rooms); err != nil {
		return nil, err
	}

	return rooms, nil
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

func (r *roomService) GetUserRole(roomID, userID int) (string, error) {
	return r.Repo.GetUserRole(roomID, userID)
}

func (s *roomService) enrichWithMemberCount(ctx context.Context, rooms []models.Room) error {
	for i, val := range rooms {
		key := fmt.Sprintf("room:%d:member_count", val.ID)
		count, err := s.Redis.Get(ctx, key).Int64()
		if err == redis.Nil {
			count = 0
		} else if err != nil {
			return fmt.Errorf("failed to get member count: %w", err)
		}
		rooms[i].MemberCount = count
	}
	return nil
}
func NewRoomService(repo repositories.RoomRepository, hub *hub.Manager,redis *redis.Client) Roomservice {
	return &roomService{Repo: repo, Hub: hub,Redis: redis}
}
