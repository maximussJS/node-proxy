package cache

import (
	"context"
	"json-rpc-node-proxy/pkg/node-proxy/cache/redis"
	"json-rpc-node-proxy/pkg/node-proxy/common/custom_errors"
	"json-rpc-node-proxy/pkg/node-proxy/config"
)

type Driver interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type Cache struct {
	driver Driver
}

func NewCache() *Cache {
	if config.CacheConfig.UseRedis == true {
		return &Cache{
			driver: redis.NewRedis(),
		}
	}

	panic("No cache driver configured")
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	select {
	case <-ctx.Done():
		return "", custom_errors.RequestTimeoutError
	default:
		return c.driver.Get(ctx, key)
	}
}

func (c *Cache) Set(ctx context.Context, key, response string) error {
	select {
	case <-ctx.Done():
		return custom_errors.RequestTimeoutError
	default:
		return c.driver.Set(ctx, key, response)
	}
}
