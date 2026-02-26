package services

import (
	"chat-app/internal/repositories"
	"chat-app/internal/shared/dtos"
	"chat-app/internal/shared/pubsub"
	"context"
	"errors"
)

type JoinRoomService interface {
	JoinRoom(ctx context.Context, userID, roomID int) (*dtos.JoinRoomResponse, error)
	LeaveRoom(ctx context.Context, userName string, roomID, userID int) error
}

type joinRoom struct {
	repo repositories.JoinRoomRepo
	pub  *pubsub.PubSub
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

	data := dtos.MapToJoinResponse(user, room)

	return data, nil
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

	return nil

}

func NewJoinRoomService(repo repositories.JoinRoomRepo, pub *pubsub.PubSub) JoinRoomService {
	return &joinRoom{
		repo: repo,
		pub:  pub,
	}
}
