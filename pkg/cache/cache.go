package cache

import (
	"context"
	"go.uber.org/dig"
	"json-rpc-node-proxy/pkg/custom_errors"
	utils_ctx "json-rpc-node-proxy/pkg/utils/ctx"
)

type IDriver interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type ICache interface {
	Get(ctx context.Context) (string, error)
	Set(ctx context.Context, data string) error
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

func (c *Cache) Get(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", custom_errors.RequestTimeoutError
	default:
		cacheKey, err := utils_ctx.GetCacheKeyFromContext(ctx)

		if err != nil {
			return "", err
		}

		return c.driver.Get(ctx, cacheKey)
	}
}

func (c *Cache) Set(ctx context.Context, data string) error {
	select {
	case <-ctx.Done():
		return custom_errors.RequestTimeoutError
	default:
		cacheKey, err := utils_ctx.GetCacheKeyFromContext(ctx)

		if err != nil {
			return err
		}

		return c.driver.Set(ctx, cacheKey, data)
	}
}
