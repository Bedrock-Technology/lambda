package main

import (
	"github.com/redis/go-redis/v9"
)

var (
	redisDB = make(map[string]*redis.Client)
)

func loadRedis(key, addr string) error {
	opt, err := redis.ParseURL(addr)
	if err != nil {
		return err
	}
	redisDB[key] = redis.NewClient(opt)
	return nil
}
