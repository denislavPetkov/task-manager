package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	InvalidKeyErr = "redis: nil"
)

type Redis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	Exists(ctx context.Context, keys ...string) (int64, error)
}

type redisInstance struct {
	redisClient *redis.Client
}

func NewRedis(config RedisConfig) Redis {
	return &redisInstance{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     config.GetAddress(),
			Username: config.GetUsername(),
			Password: config.GetPassword(),
			DB:       config.GetDbInstance(),
		}),
	}
}
