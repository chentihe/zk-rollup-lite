package cache

import (
	"context"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(context context.Context, config *config.Redis) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr(),
		Password: config.Password,
		DB:       0,
	})
	if err := client.Ping(context).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client,
	}, nil
}

func (cache *RedisCache) Get(context context.Context, key string) (string, error) {
	value, err := cache.Client.Get(context, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (cache *RedisCache) Set(context context.Context, key string, value interface{}) error {
	return cache.Client.Set(context, key, value, -1).Err()
}

func (cache *RedisCache) Close() error {
	return cache.Client.Close()
}
