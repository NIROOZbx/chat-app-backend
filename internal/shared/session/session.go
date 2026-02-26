package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const sessionTTL = 24 * time.Hour
const keyPrefix = "session:"

type Data struct {
	ID       int    `json:"ID"`
	UserName string `json:"UserName"`
}
type Store struct {
	redis *redis.Client
}

func (s *Store) CreateSession(ctx context.Context, data *Data) (string, error) {

	sessionID := uuid.NewString()

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session: %w", err)
	}
	if err := s.redis.Set(ctx, keyPrefix+sessionID, bytes, sessionTTL).Err(); err != nil {
		return "", fmt.Errorf("failed to store session: %w", err)
	}
	return sessionID, nil
}

func (s *Store) Delete(ctx context.Context, sessionID string) error {
	return s.redis.Del(ctx, keyPrefix+sessionID).Err()
}

func (s *Store) Get(ctx context.Context, sessionID string) (*Data, error) {

	val, err := s.redis.Get(ctx, keyPrefix+sessionID).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	var data Data
	err = json.Unmarshal([]byte(val), &data)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)

	}

	s.redis.Expire(ctx, keyPrefix+sessionID, sessionTTL)

	return &data, nil
}

func CreateStore(r *redis.Client) *Store {
	return &Store{redis: r}
}
