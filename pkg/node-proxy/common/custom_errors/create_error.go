package custom_errors

import (
	"fmt"
	"json-rpc-node-proxy/pkg/node-proxy/json-rpc-types"
)

var (
	cacheDriverGetErrorCode                = -32000
	cacheDriverSetErrorCode                = -32001
	keyGenerationErrorCode                 = -32002
	nodeRequestJsonMarshalErrorCode        = -32003
	nodeRequestNewRequestErrorCode         = -32004
	nodeRequestClientDoErrorCode           = -32005
	nodeRequestReadResponseBodyErrorCode   = -32006
	nodeResponseResultMarshalErrorCode     = -32007
	proxyFromCacheResultUnmarshalErrorCode = -32008
	nodeResultMarshalErrorCode             = -32009
	cacheDriverSetExpireErrorCode          = -32010
)

func NewValidationError(err error) *json_rpc_types.RpcError {
	return json_rpc_types.NewRpcError(-32602, err.Error(), nil)
}

func createJsonRpcRequestProcessingError(code int) *json_rpc_types.RpcError {
	return json_rpc_types.NewRpcError(code, "Internal Request Processing Error. Provide code to support team", nil)
}

func CreateJsonRpcError(err error) *json_rpc_types.RpcError {
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
	case ProxyFromCacheResultUnmarshalError:
		return createJsonRpcRequestProcessingError(proxyFromCacheResultUnmarshalErrorCode)
	case NodeResultMarshalError:
		return createJsonRpcRequestProcessingError(nodeResultMarshalErrorCode)
	case CacheDriverSetExpireError:
		return createJsonRpcRequestProcessingError(cacheDriverSetExpireErrorCode)
	default:
		panic(fmt.Errorf("CreateJsonRpcError Unknown error type %v", err))
	}
}
