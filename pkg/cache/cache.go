package cache

import (
	"context"
	"go.uber.org/dig"
	"json-rpc-node-proxy/internal/common/custom_errors"
)

type IDriver interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type ICache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, data string) error
}

type Cache struct {
	driver IDriver
}

type CacheDependencies struct {
	dig.In

	Driver IDriver `name:"CacheDriver"`
}

func NewCache(deps CacheDependencies) *Cache {
	return &Cache{
		driver: deps.Driver,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	select {
	case <-ctx.Done():
		return "", custom_errors.RequestTimeoutError
	default:
		return c.driver.Get(ctx, key)
	}
}

func (c *Cache) Set(ctx context.Context, key, data string) error {
	select {
	case <-ctx.Done():
		return custom_errors.RequestTimeoutError
	default:
		return c.driver.Set(ctx, key, data)
	}
}
