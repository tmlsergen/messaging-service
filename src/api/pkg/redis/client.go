package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisNilError = redis.Nil

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(dsn string) *RedisClient {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Get(c context.Context, key string) (string, error) {
	return r.client.Get(c, key).Result()
}

func (r *RedisClient) Set(c context.Context, key string, value string) error {
	return r.client.Set(c, key, value, 0).Err()
}
