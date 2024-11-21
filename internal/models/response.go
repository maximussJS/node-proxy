package models

import (
	"bytes"
	"encoding/json"
)

type JsonRpcResponse struct {
	JsonRpc string      `json:"jsonrpc"`          // The version of the JSON-RPC protocol.
	Id      interface{} `json:"id"`               // The ID of the request (can be string, number, or null).
	Result  interface{} `json:"result,omitempty"` // The result of the request (present only if successful).
	Error   *RpcError   `json:"error,omitempty"`  // The error object (present only if an error occurred).
}

func NewJsonRpcResponse(id interface{}, jsonRpc string, result interface{}, err *RpcError) *JsonRpcResponse {
	return &JsonRpcResponse{
		JsonRpc: jsonRpc,
		Id:      id,
		Result:  result,
		Error:   err,
	}
}

func NewJsonRpcResponseFromString(str string) *JsonRpcResponse {
	var response JsonRpcResponse

	if err := json.NewDecoder(bytes.NewReader([]byte(str))).Decode(&response); err != nil {
		panic(err)
	}

	return &response
}

func (r *JsonRpcResponse) IsNotError() bool {
	return r.Error == nil
}

func (r *JsonRpcResponse) CopyWithNewId(id interface{}) *JsonRpcResponse {
	return &JsonRpcResponse{
		JsonRpc: r.JsonRpc,
		Id:      id,
		Result:  r.Result,
		Error:   r.Error,
	}
}
