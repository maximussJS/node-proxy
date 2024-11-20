package config

import "time"

type IHttpServerConfig interface {
	GetPort() string
	GetJsonRpcEndpoint() string
	GetMaxRequestBodySize() int64
	GetTimeout() time.Duration
}

type HttpServerConfig struct {
	port               string
	jsonRpcEndpoint    string
	maxRequestBodySize int64
	timeout            int
}

func (hsc *HttpServerConfig) GetPort() string {
	return hsc.port
}

func (hsc *HttpServerConfig) GetJsonRpcEndpoint() string {
	return hsc.jsonRpcEndpoint
}

func (hsc *HttpServerConfig) GetMaxRequestBodySize() int64 {
	return hsc.maxRequestBodySize
}

func (hsc *HttpServerConfig) GetTimeout() time.Duration {
	return time.Duration(hsc.timeout) * time.Second
}

var httpServerConfig *HttpServerConfig

func init() {
	initHttpServerConfig()
}

func initHttpServerConfig() {
	port := EnvOptionalString("HTTP_SERVER_PORT", ":8080")

	endpoint := EnvOptionalString("HTTP_SERVER_ENDPOINT", "/")

	timeout := EnvOptionalInt("HTTP_SERVER_TIMEOUT", 30)

	maxRequestBodySize := EnvOptionalInt64("HTTP_SERVER_MAX_REQUEST_BODY_SIZE", 1024*1024)

	httpServerConfig = &HttpServerConfig{
		port:               port,
		jsonRpcEndpoint:    endpoint,
		maxRequestBodySize: maxRequestBodySize,
		timeout:            timeout,
	}
}

func SingletonHttpServerConfig() *HttpServerConfig {
	if httpServerConfig == nil {
		initHttpServerConfig()
	}

	return httpServerConfig
}
