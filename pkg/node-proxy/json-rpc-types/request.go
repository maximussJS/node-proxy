package json_rpc_types

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type JsonRpcRequest struct {
	Jsonrpc string        `json:"jsonrpc" validate:"required,eq=2.0"` // Must be "2.0"
	Method  string        `json:"method" validate:"required"`         // Method is required
	Params  []interface{} `json:"params"`                             // Params can be any array
	Id      string        `json:"id" validate:"required"`             // ID is required
}

func (r *JsonRpcRequest) Validate() error {
	validate := validator.New()

	if err := validate.Struct(r); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Error())
		}
	}

	if r.Params == nil {
		return errors.New("params must not be nil")
	}

	return nil
}

func (r *JsonRpcRequest) ErrorResponse(err *RpcError) *JsonRpcResponse {
	return NewJsonRpcResponse(r.Id, r.Jsonrpc, nil, err)
}

func (r *JsonRpcRequest) SuccessResponse(result interface{}) *JsonRpcResponse {
	return NewJsonRpcResponse(r.Id, r.Jsonrpc, result, nil)
}

func (r *JsonRpcRequest) CopyWithNewId(id string) *JsonRpcRequest {
	return &JsonRpcRequest{
		Jsonrpc: r.Jsonrpc,
		Method:  r.Method,
		Params:  r.Params,
		Id:      id,
	}
}

func (r *JsonRpcRequest) IsCacheable() bool {
	return true
}