package conn

import (
	"github.com/redis/go-redis/v9"
	"github.com/whotterre/tiermaster/internal/config"
)

// Initialize Redis client
func GetRedisClient(cfg config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	return client, nil
}