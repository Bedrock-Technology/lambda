package core

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func RedisHGet(r *redis.Client, key, field string) (string, error) {
	if r == nil {
		return "", ErrNilDB
	}

	got, err := r.HGet(context.Background(), key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return got, nil
}

func RedisHSet(r *redis.Client, key string, values ...any) (int64, error) {
	if r == nil {
		return 0, ErrNilDB
	}

	return r.HSet(context.Background(), key, values...).Result()
}
