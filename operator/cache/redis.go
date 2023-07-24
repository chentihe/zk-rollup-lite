package cache

import (
	"context"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/metrics"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client  *redis.Client
	marshal *marshaler.Marshaler
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

	redisInstance := redis_store.NewRedis(client)

	redisCacheManager := cache.New[any](redisInstance)

	promMetrics := metrics.NewPrometheus(config.Prometheus)

	cacheManager := cache.NewMetric[any](
		promMetrics,
		redisCacheManager,
	)

	marshal := marshaler.New(cacheManager)

	return &RedisCache{
		client,
		marshal,
	}, nil
}

func (cache *RedisCache) Get(context context.Context, key string, valueStruct interface{}) (interface{}, error) {
	object, err := cache.marshal.Get(context, key, valueStruct)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (cache *RedisCache) Set(context context.Context, key string, value interface{}) error {
	return cache.marshal.Set(context, key, value)
}

func (cache *RedisCache) Close() error {
	return cache.Client.Close()
}
