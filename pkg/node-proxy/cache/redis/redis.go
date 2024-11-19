package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"json-rpc-node-proxy/pkg/node-proxy/common/custom_errors"
	"json-rpc-node-proxy/pkg/node-proxy/config"
	"log"
	"time"
)

type Redis struct {
	client     *redis.Client
	defaultTTl time.Duration
}

func NewRedis() *Redis {
	return &Redis{
		client:     createClient(),
		defaultTTl: time.Duration(config.RedisConfig.DefaultTTl) * time.Second,
	}
}

func (r *Redis) Set(ctx context.Context, key string, data string) error {
	err := r.client.Set(ctx, key, data, r.defaultTTl).Err()

	if err != nil {
		log.Printf("Redis.Set() error %v", err)
		return custom_errors.CacheDriverSetError
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		log.Printf("Redis.Get() error %v", err)
		return "", custom_errors.CacheDriverGetError
	}

	if _, err := r.client.Expire(ctx, key, r.defaultTTl).Result(); err != nil {
		log.Printf("Redis.Expire() error %v", err)
		return "", custom_errors.CacheDriverSetError
	}

	return val, nil
}

func createClient() *redis.Client {
	opts, err := redis.ParseURL(config.RedisConfig.Url)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opts)
}
