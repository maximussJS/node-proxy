package models

import "fmt"

type RpcError struct {
	Code    int    `json:"code"`    // The error code.
	Message string `json:"message"` // A short description of the error.
	Data    any    `json:"data"`    // Additional information about the error (optional).
}

func NewRpcError(code int, message string, data any) *RpcError {
	return &RpcError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (e *RpcError) Error() string {
	return fmt.Sprintf("Error code: %d, message: %s, data: %v", e.Code, e.Message, e.Data)
}
