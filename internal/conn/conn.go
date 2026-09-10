package conn

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/whotterre/tiermaster/internal/config"
)

// GetRedisClient initializes and returns a Redis client with a retry mechanism.
func GetRedisClient(cfg config.Config) (*redis.Client, error) {
	var client *redis.Client

	connectEffector := func(ctx context.Context) (string, error) {
		opt, err := redis.ParseURL(cfg.RedisAddr)
	    if err != nil {
	        return "", err
	    }

		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
		}

		 c := redis.NewClient(opt)

		_, err = c.Ping().Result()
	    if err != nil {
	        return "", err
	    }

   		client = c
		return "Redis connected", nil
	}

	// Retry connecting to Redis 5 times with a 2-second delay between attempts
	retryConnect := Retry(connectEffector, 5, 2*time.Second)

	_, err := retryConnect(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}

type Effector func(context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				return response, err
			}

			log.Printf("Attempt %d failed, retrying in %v", r+1, delay)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}
