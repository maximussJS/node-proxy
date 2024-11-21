package custom_errors

import (
	"fmt"
	"json-rpc-node-proxy/internal/models"
)

var (
	cacheDriverGetErrorCode              = -32000
	cacheDriverSetErrorCode              = -32001
	keyGenerationErrorCode               = -32002
	nodeRequestJsonMarshalErrorCode      = -32003
	nodeRequestNewRequestErrorCode       = -32004
	nodeRequestClientDoErrorCode         = -32005
	nodeRequestReadResponseBodyErrorCode = -32006
	nodeResponseResultMarshalErrorCode   = -32007
	cacheResultUnmarshalErrorCode        = -32008
	nodeResultMarshalErrorCode           = -32009
	cacheDriverSetExpireErrorCode        = -32010
	ctxCacheKeyEmptyErrorCode            = -32011
	ctxJsonRpcRequestEmptyError          = -32012
)

func NewValidationError(err error) *models.RpcError {
	return models.NewRpcError(-32602, err.Error(), nil)
}

func createJsonRpcRequestProcessingError(code int) *models.RpcError {
	return models.NewRpcError(code, "Internal Request Processing Error. Provide code to support team", nil)
}

func CreateJsonRpcError(err error) *models.RpcError {
	switch err {
	case CacheDriverGetError:
		return createJsonRpcRequestProcessingError(cacheDriverGetErrorCode)
	case CacheDriverSetError:
		return createJsonRpcRequestProcessingError(cacheDriverSetErrorCode)
	case KeyGenerationError:
		return createJsonRpcRequestProcessingError(keyGenerationErrorCode)
	case NodeRequestJsonMarshalError:
		return createJsonRpcRequestProcessingError(nodeRequestJsonMarshalErrorCode)
	case NodeRequestNewRequestError:
		return createJsonRpcRequestProcessingError(nodeRequestNewRequestErrorCode)
	case NodeRequestClientDoError:
		return createJsonRpcRequestProcessingError(nodeRequestClientDoErrorCode)
	case NodeRequestReadResponseBodyError:
		return createJsonRpcRequestProcessingError(nodeRequestReadResponseBodyErrorCode)
	case NodeResponseResultMarshalError:
		return createJsonRpcRequestProcessingError(nodeResponseResultMarshalErrorCode)
	case CacheResultUnmarshalError:
		return createJsonRpcRequestProcessingError(cacheResultUnmarshalErrorCode)
	case NodeResultMarshalError:
		return createJsonRpcRequestProcessingError(nodeResultMarshalErrorCode)
	case CacheDriverSetExpireError:
		return createJsonRpcRequestProcessingError(cacheDriverSetExpireErrorCode)
	case CtxCacheKeyEmptyError:
		return createJsonRpcRequestProcessingError(ctxCacheKeyEmptyErrorCode)
	case CtxJsonRpcRequestEmptyError:
		return createJsonRpcRequestProcessingError(ctxJsonRpcRequestEmptyError)

	default:
		panic(fmt.Errorf("CreateJsonRpcError Unknown error type %v", err))
	}
}
