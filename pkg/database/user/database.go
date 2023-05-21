package user

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	InvalidKeyErr = "redis: nil"
)

var (
	logger *zap.Logger
)

func init() {
	logger = zap.L().Named("user-db")
}

type Redis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Del(ctx context.Context, keys ...string) (int64, error)
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
