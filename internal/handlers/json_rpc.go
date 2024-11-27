package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/dig"
	"io"
	"json-rpc-node-proxy/internal/models"
	"json-rpc-node-proxy/internal/services"
	"json-rpc-node-proxy/pkg/config"
	"json-rpc-node-proxy/pkg/custom_errors"
	"json-rpc-node-proxy/pkg/logger"
	"json-rpc-node-proxy/pkg/utils/responses"
	"net/http"
	"time"
)

type IJsonRpcHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type JsonRpcHandler struct {
	logger             logger.ILogger
	requestTimeout     time.Duration
	maxRequestBodySize int64
	jsonRpc            services.IJsonRpcService
}

type JsonRpcRequestHandlerDependencies struct {
	dig.In
	Logger  logger.ILogger           `name:"Logger"`
	Cfg     config.IConfig           `name:"Config"`
	JsonRpc services.IJsonRpcService `name:"JsonRpcService"`
}

func NewJsonRpcHandler(deps JsonRpcRequestHandlerDependencies) *JsonRpcHandler {
	return &JsonRpcHandler{
		logger:             deps.Logger,
		requestTimeout:     deps.Cfg.GetHttpRequestTimeout(),
		maxRequestBodySize: deps.Cfg.GetHttpMaxRequestBodySize(),
		jsonRpc:            deps.JsonRpc,
	}
}

func (h *JsonRpcHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	h.processJsonRpcRequest(w, r.WithContext(ctx))
}

func (h *JsonRpcHandler) processJsonRpcRequest(w http.ResponseWriter, r *http.Request) {
	request, err := h.parseJsonRpcRequest(w, r)

	if err != nil {
		responses.BadRequest(w, fmt.Errorf(`{"error": "Error unmarshalling JSON: %v"}`, err.Error()))
		return
	}

	validationErrorResponseBytes, err := h.validateJsonRpcRequest(request)

	if err != nil {
		responses.InternalServerError(w)
		return
	}

	if validationErrorResponseBytes != nil {
		responses.Success(w, validationErrorResponseBytes)
		return
	}

	select {
	case <-r.Context().Done():
		responses.RequestTimeout(w)
		return
	default:
		response, err := h.jsonRpc.HandleRequest(r.Context(), request)

		if err != nil {
			if errors.Is(err, custom_errors.RequestTimeoutError) {
				responses.RequestTimeout(w)
				return
			}

			serr, ok := err.(*models.RpcError)
			if ok {
				responseBytes, err := json.Marshal(request.ErrorResponse(serr))

				if err != nil {
					h.logger.Error(fmt.Sprintf("Error marshalling error response: %v", err))
					responses.InternalServerError(w)
					return
				}

				responses.Success(w, responseBytes)
				return
			}

			h.logger.Error(fmt.Sprintf("Unhandled error while processing request: %v", err))

			responses.InternalServerError(w)
			return
		}

		bytes, err := json.Marshal(response)

		if err != nil {
			h.logger.Error(fmt.Sprintf("Error marshalling response: %v", err))
			responses.InternalServerError(w)
			return
		}

		responses.Success(w, bytes)
	}
}

func (h *JsonRpcHandler) parseJsonRpcRequest(w http.ResponseWriter, r *http.Request) (*models.JsonRpcRequest, error) {
	var request models.JsonRpcRequest

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, h.maxRequestBodySize))

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error(fmt.Sprintf("Error closing request body: %v", err))
		}
	}(r.Body)

	err := decoder.Decode(&request)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (h *JsonRpcHandler) validateJsonRpcRequest(request *models.JsonRpcRequest) ([]byte, error) {
	if err := request.Validate(); err != nil {
		validationErrorResponseBytes, err := json.Marshal(request.ErrorResponse(custom_errors.NewValidationError(err)))

		if err != nil {
			h.logger.Error(fmt.Sprintf("Error marshalling error response: %v", err))
			return nil, err
		}

		return validationErrorResponseBytes, nil
	}

	return nil, nil
}
