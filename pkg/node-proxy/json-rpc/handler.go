package json_rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"json-rpc-node-proxy/pkg/node-proxy/common/custom_errors"
	"json-rpc-node-proxy/pkg/node-proxy/json-rpc-types"
	"json-rpc-node-proxy/pkg/node-proxy/utils/responses"
	"json-rpc-node-proxy/pkg/node-proxy/worker_pool"
	"log"
	"net/http"
	"time"
)

type JsonRpcHandler struct {
	requestTimeout     time.Duration
	maxRequestBodySize int64
	pool               *worker_pool.WorkerPool[*json_rpc_types.JsonRpcResponse]
	proxy              *Proxy
}

func NewJsonRpcHandler(maxRequestBodySize int64, timeout time.Duration) *JsonRpcHandler {
	return &JsonRpcHandler{
		requestTimeout:     timeout,
		maxRequestBodySize: maxRequestBodySize,
		pool:               worker_pool.NewWorkerPool[*json_rpc_types.JsonRpcResponse](),
		proxy:              NewProxy(),
	}
}

func (h *JsonRpcHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	h.pool.Submit(func() (*json_rpc_types.JsonRpcResponse, error) {
		select {
		case <-ctx.Done():
			responses.RequestTimeout(w)
		default:
			h.processRequest(w, r.WithContext(ctx))
		}

		return nil, nil // This is a dummy return value, it's not used, but it's required by the pond.WorkerPool
	})
}

func (h *JsonRpcHandler) processRequest(w http.ResponseWriter, r *http.Request) {
	var request json_rpc_types.JsonRpcRequest

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, h.maxRequestBodySize))

	defer r.Body.Close()

	err := decoder.Decode(&request)

	if err != nil {
		responses.BadRequest(w, fmt.Errorf(`{"error": "Error unmarshalling JSON: %v"}`, err.Error()))
		return
	}

	if err := request.Validate(); err != nil {
		responseBytes, err := json.Marshal(request.ErrorResponse(custom_errors.NewValidationError(err)))

		if err != nil {
			log.Printf("Error marshalling error response: %v", err)
			responses.InternalServerError(w)
			return
		}

		responses.Success(w, responseBytes)
		return
	}

	select {
	case <-r.Context().Done():
		responses.RequestTimeout(w)
		return
	default:
		response, err := h.proxy.HandleRequest(r.Context(), &request)

		if err != nil {
			if errors.Is(err, custom_errors.RequestTimeoutError) {
				responses.RequestTimeout(w)
				return
			}

			serr, ok := err.(*json_rpc_types.RpcError)
			if ok {
				responseBytes, err := json.Marshal(request.ErrorResponse(serr))

				if err != nil {
					log.Printf("Error marshalling error response: %v", err)
					responses.InternalServerError(w)
					return
				}

				responses.Success(w, responseBytes)
				return
			}

			log.Printf("Unhandled error while processing request: %v", err)

			responses.InternalServerError(w)
			return
		}

		bytes, err := json.Marshal(response)

		if err != nil {
			log.Printf("Error marshalling response: %v", err)
			responses.InternalServerError(w)
			return
		}

		responses.Success(w, bytes)
	}
}
