package config

type httpServerConfig struct {
	Port               string
	JsonRpcEndpoint    string
	MaxRequestBodySize int64
	Timeout            int
}

var HttpServerConfig *httpServerConfig

func init() {
	port := EnvOptionalString("HTTP_SERVER_PORT", ":8080")

	endpoint := EnvOptionalString("HTTP_SERVER_ENDPOINT", "/")

	timeout := EnvOptionalInt("HTTP_SERVER_TIMEOUT", 30)

	maxRequestBodySize := EnvOptionalInt64("HTTP_SERVER_MAX_REQUEST_BODY_SIZE", 1024*1024)

	HttpServerConfig = &httpServerConfig{
		Port:               port,
		JsonRpcEndpoint:    endpoint,
		MaxRequestBodySize: maxRequestBodySize,
		Timeout:            timeout,
	}
}
