package ctx

import (
	"context"
	"json-rpc-node-proxy/internal/models"
	"json-rpc-node-proxy/pkg/custom_errors"
)

var jsonRpcRequestKey = "jsonRpcRequest"

func GetJsonRpcRequestFromContext(ctx context.Context) (*models.JsonRpcRequest, error) {
	if result, ok := ctx.Value(jsonRpcRequestKey).(*models.JsonRpcRequest); ok {
		return result, nil
	}

	return nil, custom_errors.CtxJsonRpcRequestEmptyError
}

func GetContextWithJsonRpcRequest(ctx context.Context, value *models.JsonRpcRequest) context.Context {
	return context.WithValue(ctx, jsonRpcRequestKey, value)
}
