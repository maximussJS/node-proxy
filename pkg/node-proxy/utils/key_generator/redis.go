package key_generator

import (
	"encoding/json"
	"fmt"
	"json-rpc-node-proxy/pkg/node-proxy/common/custom_errors"
	"log"
	"strings"
)

type RedisKeyGenerator struct {
}

func NewRedisKeyGenerator() *RedisKeyGenerator {
	return &RedisKeyGenerator{}
}

func (r *RedisKeyGenerator) GenerateJsonRpcKey(jsonrpc, method string, params []any) (string, error) {
	paramStrings := make([]string, len(params))
	for i, param := range params {
		paramBytes, err := json.Marshal(param)

		if err != nil {
			log.Printf("RedisKeyGenerator.GenerateKey() error %v", err)
			return "", custom_errors.KeyGenerationError
		}

		paramStrings[i] = string(paramBytes)
	}

	paramsPart := strings.Join(paramStrings, ",")

	return fmt.Sprintf("request-key:%s:%s:%s", jsonrpc, method, paramsPart), nil
}
