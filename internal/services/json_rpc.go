package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/dig"
	"json-rpc-node-proxy/internal/models"
	"json-rpc-node-proxy/pkg/cache"
	"json-rpc-node-proxy/pkg/custom_errors"
	"json-rpc-node-proxy/pkg/key_generator"
	"json-rpc-node-proxy/pkg/logger"
	utils_ctx "json-rpc-node-proxy/pkg/utils/ctx"
)

type IJsonRpcService interface {
	HandleRequest(ctx context.Context, request *models.JsonRpcRequest) (*models.JsonRpcResponse, error)
}

type JsonRpcService struct {
	logger       logger.ILogger
	node         INodeService
	cache        cache.ICache
	keyGenerator key_generator.IKeyGenerator
}

type JsonRpcServiceDependencies struct {
	dig.In

	Logger       logger.ILogger              `name:"Logger"`
	Node         INodeService                `name:"NodeService"`
	Cache        cache.ICache                `name:"Cache"`
	KeyGenerator key_generator.IKeyGenerator `name:"KeyGenerator"`
}

func NewJsonRpcService(deps JsonRpcServiceDependencies) *JsonRpcService {
	return &JsonRpcService{
		logger:       deps.Logger,
		node:         deps.Node,
		cache:        deps.Cache,
		keyGenerator: deps.KeyGenerator,
	}
}

func (p *JsonRpcService) HandleRequest(ctx context.Context, request *models.JsonRpcRequest) (*models.JsonRpcResponse, error) {
	select {
	case <-ctx.Done():
		return nil, custom_errors.RequestTimeoutError
	default:
		cacheKey, err := p.keyGenerator.GenerateJsonRpcKey(request.Jsonrpc, request.Method, request.Params)

		if err != nil {
			return nil, custom_errors.CreateJsonRpcError(err)
		}

		ctxWithValues := utils_ctx.GetContextWithJsonRpcRequest(utils_ctx.GetContextWithCacheKey(ctx, cacheKey), request)

		fromCache, err := p.tryToFindInCache(ctxWithValues)

		if err != nil {
			return nil, custom_errors.CreateJsonRpcError(err)
		}

		if fromCache != nil {
			return fromCache, nil
		}

		response, err := p.node.Request(ctxWithValues)

		if err != nil {
			return nil, custom_errors.CreateJsonRpcError(err)
		}

		if err := p.saveToCacheIfNeed(ctxWithValues, response); err != nil {
			return nil, custom_errors.CreateJsonRpcError(err)
		}

		return response, nil
	}
}

func (p *JsonRpcService) tryToFindInCache(ctx context.Context) (*models.JsonRpcResponse, error) {
	fromCache, err := p.cache.Get(ctx)

	if err != nil {
		return nil, err
	}

	request, err := utils_ctx.GetJsonRpcRequestFromContext(ctx)

	if err != nil {
		return nil, err
	}

	if fromCache != "" {
		p.logger.Debug(fmt.Sprintf("Response for request %s was found in cache", request))

		var fromCacheResult interface{}

		if err := json.Unmarshal([]byte(fromCache), &fromCacheResult); err != nil {
			p.logger.Error(fmt.Sprintf("Failed to unmarshal from cache: %v", err))

			return nil, custom_errors.CacheResultUnmarshalError
		}

		return request.SuccessResponse(fromCacheResult), nil
	}

	return nil, nil
}

func (p *JsonRpcService) saveToCacheIfNeed(ctx context.Context, response *models.JsonRpcResponse) error {
	request, err := utils_ctx.GetJsonRpcRequestFromContext(ctx)

	if err != nil {
		return err
	}

	if request.IsCacheable() && response.IsNotError() {
		result, err := json.Marshal(response.Result)

		if err != nil {
			return custom_errors.NodeResultMarshalError
		}

		if err := p.cache.Set(ctx, string(result)); err != nil {
			return err
		}

		p.logger.Debug(fmt.Sprintf("Result for request %s was saved in cache", request))
	} else {
		p.logger.Debug(fmt.Sprintf("Result for request %s was not saved in cache", request))
	}

	return nil
}
