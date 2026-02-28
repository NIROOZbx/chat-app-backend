package services

import (
	"chat-app/internal/repositories"
	"chat-app/internal/shared/dtos"
	"chat-app/internal/shared/pubsub"
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type JoinRoomService interface {
	JoinRoom(ctx context.Context, userID, roomID int) (*dtos.JoinRoomResponse, error)
	JoinPrivateRoom(ctx context.Context, userID int, inviteCode string) (*dtos.JoinRoomResponse, error)
	LeaveRoom(ctx context.Context, userName string, roomID, userID int) error
}

type joinRoom struct {
	repo  repositories.JoinRoomRepo
	pub   *pubsub.PubSub
	redis *redis.Client
}

func (s *joinRoom) JoinRoom(ctx context.Context, userID, roomID int) (*dtos.JoinRoomResponse, error) {
	isMember, err := s.repo.IsRoomMember(userID, roomID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, errors.New("user is already a member of this room")
	}
	room, err := s.repo.GetRoomByID(roomID)
	if err != nil {
		return nil, err
	}

	if room.IsPrivate {
		return nil, errors.New("cannot join a private room without an invite code")
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.JoinRoom(userID, roomID, "member"); err != nil {
		return nil, err
	}
	role, err := s.repo.GetUserRole(roomID, userID)
	if err != nil {
		return nil, err
	}

	event := &pubsub.UserJoinedEvent{
		BaseEvent: pubsub.BaseEvent{
			Type: pubsub.EventUserJoined,
		},
		RoomID:   room.ID,
		UserID:   user.ID,
		UserName: user.UserName,
		Role:     role,
	}

	err = s.pub.Publish(ctx, roomID, event)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("room:%d:member_count", roomID)
	s.redis.Incr(ctx, key)

	data := dtos.MapToJoinResponse(user, room)

	return data, nil
}

func (s *joinRoom) JoinPrivateRoom(ctx context.Context, userID int, inviteCode string) (*dtos.JoinRoomResponse, error) {
	room, err := s.repo.GetRoomByInviteCode(inviteCode)
	if err != nil {
		return nil, errors.New("invalid or expired invite code")
	}

	return s.JoinRoom(ctx, userID, room.ID)
}

func (s *joinRoom) LeaveRoom(ctx context.Context, userName string, roomID, userID int) error {

	err := s.repo.LeaveRoom(roomID, userID)

	if err != nil {
		return err
	}
	payload := &pubsub.UserLeaveEvent{
		BaseEvent: pubsub.BaseEvent{
			Type: pubsub.EventUserLeft,
		},
		RoomID:   roomID,
		UserID:   userID,
		UserName: userName,
	}

	err = s.pub.Publish(ctx, roomID, payload)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("room:%d:member_count", roomID)
	s.redis.Decr(ctx, key)

	return nil

}

func NewJoinRoomService(repo repositories.JoinRoomRepo, pub *pubsub.PubSub, redis *redis.Client) JoinRoomService {
	return &joinRoom{
		repo:  repo,
		pub:   pub,
		redis: redis,
	}
}
