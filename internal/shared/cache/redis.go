package cache

import (
	"chat-app/internal/shared/config"
	"log"

	"github.com/redis/go-redis/v9"
)

func SetupRedis(cfg *config.RedisConfig) *redis.Client {

	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {	
		log.Fatalf("Error parsing redis url: %v", err)
	}

	client := redis.NewClient(opt)

	log.Println("✅ Redis connected")

	return client
}
