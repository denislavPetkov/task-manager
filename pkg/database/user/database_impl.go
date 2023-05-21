package user

import (
	"context"
	"time"
)

func (r *redisInstance) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.redisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logger.Error("Failed to set a new key")
		return err
	}

	logger.Info("Set a new key successfully")

	return err
}

func (r *redisInstance) Get(ctx context.Context, key string) (string, error) {
	result, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		logger.Error("Failed to get a key")
		return result, err
	}

	logger.Info("Got a key successfully")

	return result, err
}

func (r *redisInstance) Exists(ctx context.Context, keys ...string) (int64, error) {
	result, err := r.redisClient.Exists(ctx, keys...).Result()
	if err != nil {
		logger.Error("Failed to check if a key exists")
		return result, err
	}

	logger.Info("Checked if a key exists successfully")

	return result, err
}

func (r *redisInstance) Del(ctx context.Context, keys ...string) (int64, error) {
	result, err := r.redisClient.Del(ctx, keys...).Result()
	if err != nil {
		logger.Error("Failed to delete keys")
		return result, err
	}

	logger.Info("Deleted keys successfully")

	return result, err
}
