package services

import (
	"chat-app/internal/models"
	"chat-app/internal/repositories"
	"chat-app/internal/shared/dtos"
	"chat-app/internal/shared/pubsub"
	"context"
	"encoding/json"

	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type MessageService interface {
	SendMessage(ctx context.Context, roomID, userID int, userName, content string) (*dtos.MessageResponse, error)
	GetMessages(ctx context.Context, roomID, limit, offset int) ([]dtos.MessageResponse, error)
}

type messageService struct {
	repo   repositories.MessageRepo
	pubsub *pubsub.PubSub
	redis  *redis.Client
}

func NewMessageService(repo repositories.MessageRepo, ps *pubsub.PubSub, redis *redis.Client) MessageService {
	return &messageService{repo: repo, pubsub: ps, redis: redis}
}

func (s *messageService) SendMessage(ctx context.Context, roomID, userID int, userName, content string) (*dtos.MessageResponse, error) {
	msg := &models.Message{
		RoomID:  roomID,
		UserID:  userID,
		Content: content,
	}

	saved, err := s.repo.Save(msg)
	if err != nil {
		return nil, err
	}

	resp := &dtos.MessageResponse{
		ID:        saved.ID,
		RoomID:    saved.RoomID,
		UserID:    saved.UserID,
		UserName:  userName,
		Content:   saved.Content,
		CreatedAt: saved.CreatedAt,
	}

	s.pubsub.Publish(ctx, roomID, pubsub.MessageEvent{
		BaseEvent: pubsub.BaseEvent{Type: pubsub.EventNewMessage},
		RoomID:    saved.RoomID,
		UserID:    saved.UserID,
		UserName:  userName,
		Content:   saved.Content,
		MessageID: saved.ID,
		SentAt:    saved.CreatedAt,
	})
	
	pattern := "room:" + strconv.Itoa(roomID) + ":p:*"
	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		s.redis.Del(ctx, iter.Val())
	}

	return resp, nil
}

func (s *messageService) GetMessages(ctx context.Context, roomID, limit, page int) ([]dtos.MessageResponse, error) {
	redisKey := "room:" + strconv.Itoa(roomID) + ":p:" + strconv.Itoa(page) + ":l:" + strconv.Itoa(limit)
	bytes, err := s.redis.Get(ctx, redisKey).Result()

	if err == redis.Nil {
		var cachedData []dtos.MessageResponse
		if unmarshalErr := json.Unmarshal([]byte(bytes), &cachedData); unmarshalErr == nil {
			return cachedData, nil
		}

	}

	offset := (page - 1) * limit
	messages, err := s.repo.GetByRoomID(roomID, limit, offset)
	if err != nil {
		return nil, err
	}

	var resp []dtos.MessageResponse
	for _, m := range messages {
		resp = append(resp, dtos.MessageResponse{
			ID:        m.ID,
			RoomID:    m.RoomID,
			UserID:    m.UserID,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
			UserName:  m.Username,
		})
	}

	if jsonData, marshalErr := json.Marshal(resp); marshalErr == nil {
		ttl := 60 * time.Minute
		s.redis.Set(ctx, redisKey, jsonData, ttl)
	}

	return resp, nil
}
