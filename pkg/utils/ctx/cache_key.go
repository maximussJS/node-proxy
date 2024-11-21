package ctx

import (
	"context"
	"json-rpc-node-proxy/pkg/custom_errors"
)

var key = "cacheKey"

func GetCacheKeyFromContext(ctx context.Context) (string, error) {
	if cacheKey, ok := ctx.Value(key).(string); ok {
		return cacheKey, nil
	}

	return "", custom_errors.CtxJsonRpcRequestEmptyError
}

func GetContextWithCacheKey(ctx context.Context, cacheKey string) context.Context {
	return context.WithValue(ctx, key, cacheKey)
}
