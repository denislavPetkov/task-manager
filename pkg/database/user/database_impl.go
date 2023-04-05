package database

import (
	"context"
	"time"
)

func (r *redisInstance) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redisClient.Set(ctx, key, value, expiration).Err()
}

func (r *redisInstance) Get(ctx context.Context, key string) (string, error) {
	return r.redisClient.Get(ctx, key).Result()
}

func (r *redisInstance) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.redisClient.HSet(ctx, key, values...).Err()
}

func (r *redisInstance) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.redisClient.Exists(ctx, keys...).Result()
}
