package redisdb

import (
	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	*redis.Client
}

func NewRedisStore() *RedisStore {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	rs := RedisStore{
		Client: c,
	}

	return &rs
}
