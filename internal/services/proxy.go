package services

import (
	"context"
	"encoding/json"
	"go.uber.org/dig"
	"json-rpc-node-proxy/internal/common/custom_errors"
	"json-rpc-node-proxy/internal/models"
	"json-rpc-node-proxy/pkg/cache"
	"json-rpc-node-proxy/pkg/key_generator"
	"log"
)

type IProxyService interface {
	HandleRequest(ctx context.Context, request *models.JsonRpcRequest) (*models.JsonRpcResponse, error)
}

type ProxyService struct {
	node         INodeService
	cache        cache.ICache
	keyGenerator key_generator.IKeyGenerator
}

type ProxyServiceDependencies struct {
	dig.In
	Node         INodeService                `name:"NodeService"`
	Cache        cache.ICache                `name:"Cache"`
	KeyGenerator key_generator.IKeyGenerator `name:"KeyGenerator"`
}

func NewProxyService(deps ProxyServiceDependencies) *ProxyService {
	return &ProxyService{
		node:         deps.Node,
		cache:        deps.Cache,
		keyGenerator: deps.KeyGenerator,
	}
}

func (p *ProxyService) HandleRequest(ctx context.Context, request *models.JsonRpcRequest) (*models.JsonRpcResponse, error) {
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
