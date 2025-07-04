package core

import (
	"context"
	"time"

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

func RedisHExpire(r *redis.Client, key string, duration string, fields ...string) ([]int64, error) {
	if r == nil {
		return nil, ErrNilDB
	}

	du, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	return r.HExpire(context.Background(), key, du, fields...).Result()
}

func RedisHKeys(r *redis.Client, key string) ([]string, error) {
	if r == nil {
		return nil, ErrNilDB
	}

	keys, err := r.HKeys(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
