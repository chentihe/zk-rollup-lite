package cache

import (
	"context"
	"os"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(context context.Context, config *config.Redis) (*RedisCache, error) {
	// for docker image to retrieve redis host
	host := os.Getenv("REDIS_HOST")
	if host != "" {
		config.Host = host
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr(),
		Password: config.Password,
		DB:       0,
	})
	if err := client.Ping(context).Err(); err != nil {
		return nil, err
	}

	cache := &RedisCache{client}
	cache.initKeys(context, &config.Keys)

	return cache, nil
}

// LastInsertedKey & RollupedTxsKey should exist in redis before the creation of L2 tx
func (cache *RedisCache) initKeys(context context.Context, keys *config.Keys) {
	if _, err := cache.Get(context, keys.LastInsertedKey); err != nil {
		cache.Set(context, keys.LastInsertedKey, "0")
	}

	if _, err := cache.Get(context, keys.RollupedTxsKey); err != nil {
		cache.Set(context, keys.RollupedTxsKey, "0")
	}
}

func (cache *RedisCache) Get(context context.Context, key string) (string, error) {
	value, err := cache.Client.Get(context, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (cache *RedisCache) Set(context context.Context, key string, value interface{}) error {
	// -1 means no expiration for cache data
	return cache.Client.Set(context, key, value, -1).Err()
}

func (cache *RedisCache) Del(context context.Context, keys []string) error {
	return cache.Client.Del(context, keys...).Err()
}

func (cache *RedisCache) Close() error {
	return cache.Client.Close()
}
