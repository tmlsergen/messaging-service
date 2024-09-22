package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

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

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}
