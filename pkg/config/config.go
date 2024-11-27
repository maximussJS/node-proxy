package config

import (
	"json-rpc-node-proxy/pkg/logger"
	"time"
)

type IConfig interface {
	UseRedis() bool
	GetHttpPort() string
	GetHttpJsonRpcEndpoint() string
	GetHttpMaxRequestBodySize() int64
	GetHttpRequestTimeout() time.Duration
	GetMaxPoolWorkers() int
	GetNodeConfigPath() string
	GetRedisUrl() string
	GetRedisDefaultTTl() time.Duration
}

type Config struct {
	logger                 logger.ILogger
	useRedis               bool
	httpPort               string
	httpJsonRpcEndpoint    string
	httpMaxRequestBodySize int64
	httpRequestTimeout     time.Duration
	maxPoolWorkers         int
	nodeConfigPath         string
	redisUrl               string
	redisDefaultTTl        time.Duration
}

var config *Config

func (c *Config) UseRedis() bool {
	return c.useRedis
}

func (c *Config) GetHttpPort() string {
	return c.httpPort
}

func (c *Config) GetHttpJsonRpcEndpoint() string {
	return c.httpJsonRpcEndpoint
}

func (c *Config) GetHttpMaxRequestBodySize() int64 {
	return c.httpMaxRequestBodySize
}

func (c *Config) GetHttpRequestTimeout() time.Duration {
	return c.httpRequestTimeout
}

func (c *Config) GetMaxPoolWorkers() int {
	return c.maxPoolWorkers
}

func (c *Config) GetNodeConfigPath() string {
	return c.nodeConfigPath
}

func (c *Config) GetRedisUrl() string {
	return c.redisUrl
}

func (c *Config) GetRedisDefaultTTl() time.Duration {
	return c.redisDefaultTTl
}

func initConfig() *Config {
	useRedis := EnvOptionalBool("USE_REDIS", true)

	httpPort := EnvOptionalString("HTTP_SERVER_PORT", ":8080")

	httpJsonRpcEndpoint := EnvOptionalString("HTTP_SERVER_ENDPOINT", "/")

	httpRequestTimeout := EnvOptionalInt("HTTP_SERVER_TIMEOUT", 30)

	maxPoolWorkers := EnvOptionalInt("HTTP_SERVER_MAX_POOL_WORKERS", 1000)

	httpMaxRequestBodySize := EnvOptionalInt64("HTTP_SERVER_MAX_REQUEST_BODY_SIZE", 1024*1024)

	nodeConfigPath := EnvRequiredString("NODE_CONFIG_PATH")

	redisUrl := EnvRequiredString("REDIS_URL")

	redisDefaultTTl := EnvOptionalInt("REDIS_DEFAULT_TTL", 30)

	return &Config{
		useRedis:               useRedis,
		httpPort:               httpPort,
		httpJsonRpcEndpoint:    httpJsonRpcEndpoint,
		httpRequestTimeout:     time.Duration(httpRequestTimeout) * time.Second,
		maxPoolWorkers:         maxPoolWorkers,
		httpMaxRequestBodySize: httpMaxRequestBodySize,
		nodeConfigPath:         nodeConfigPath,
		redisUrl:               redisUrl,
		redisDefaultTTl:        time.Duration(redisDefaultTTl) * time.Second,
	}
}

func SingletonConfig() *Config {
	if config == nil {
		config = initConfig()
	}

	return config
}
