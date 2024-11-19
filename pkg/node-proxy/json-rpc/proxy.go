package json_rpc

import (
	"context"
	"encoding/json"
	"json-rpc-node-proxy/pkg/node-proxy/cache"
	"json-rpc-node-proxy/pkg/node-proxy/common/custom_errors"
	"json-rpc-node-proxy/pkg/node-proxy/config"
	"json-rpc-node-proxy/pkg/node-proxy/json-rpc-types"
	"json-rpc-node-proxy/pkg/node-proxy/node"
	"json-rpc-node-proxy/pkg/node-proxy/utils/key_generator"
	"log"
)

type Proxy struct {
	node         *node.Node
	cache        *cache.Cache
	keyGenerator key_generator.KeyGenerator
}

func NewProxy() *Proxy {
	if config.CacheConfig.UseRedis == false {
		panic("Only Redis cache is supported for now")
	}

	return &Proxy{
		node:         node.NewNode(),
		cache:        cache.NewCache(),
		keyGenerator: key_generator.NewRedisKeyGenerator(),
	}
}

func (p *Proxy) HandleRequest(ctx context.Context, request *json_rpc_types.JsonRpcRequest) (*json_rpc_types.JsonRpcResponse, error) {
	select {
	case <-ctx.Done():
		return nil, custom_errors.RequestTimeoutError
	default:
		cacheKey, err := p.keyGenerator.GenerateJsonRpcKey(request.Jsonrpc, request.Method, request.Params)

		if err != nil {
			return nil, custom_errors.CreateJsonRpcError(err)
		}

		fromCache, err := p.cache.Get(ctx, cacheKey)

		if err != nil {
			return nil, custom_errors.CreateJsonRpcError(err)
		}

		if fromCache != "" {
			var fromCacheResult interface{}

			if err := json.Unmarshal([]byte(fromCache), &fromCacheResult); err != nil {
				log.Printf("Failed to unmarshal from cache: %v", err)
				return nil, custom_errors.CreateJsonRpcError(custom_errors.ProxyFromCacheResultUnmarshalError)
			}

			return request.SuccessResponse(fromCacheResult), nil
		}

		response, err := p.node.Request(ctx, request)

		if err != nil {
			return nil, custom_errors.CreateJsonRpcError(err)
		}

		if request.IsCacheable() && response.IsCacheable() {
			result, err := json.Marshal(response.Result)

			if err != nil {
				return nil, custom_errors.CreateJsonRpcError(custom_errors.NodeResultMarshalError)
			}

			if err := p.cache.Set(ctx, cacheKey, string(result)); err != nil {
				return nil, custom_errors.CreateJsonRpcError(err)
			}
		}

		return response, nil
	}
}
