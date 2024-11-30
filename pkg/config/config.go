package config

import (
	"errors"
	"fmt"
	"go.uber.org/dig"
	"gopkg.in/yaml.v3"
	"json-rpc-node-proxy/pkg/logger"
	"os"
	"time"
)

type ConfigDeps struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`
}

type httpConfig struct {
	Port               string `yaml:"port"`
	JsonRpcEndpoint    string `yaml:"jsonRpcEndpoint"`
	MaxRequestBodySize int64  `yaml:"maxRequestBodySize"`
	RequestTimeout     int    `yaml:"requestTimeout"`
	MaxPoolWorkers     int    `yaml:"maxPoolWorkers"`
}

func (h *httpConfig) String() string {
	return fmt.Sprintf("HttpConfig{Port: %s, JsonRpcEndpoint: %s, MaxRequestBodySize: %d, RequestTimeout: %d, MaxPoolWorkers: %d}",
		h.Port, h.JsonRpcEndpoint, h.MaxRequestBodySize, h.RequestTimeout, h.MaxPoolWorkers)
}

type redisConfig struct {
	Url        string `yaml:"url"`
	DefaultTTL int    `yaml:"defaultTTL"`
}

func (r *redisConfig) String() string {
	return fmt.Sprintf("RedisConfig{Url: %s, DefaultTTL: %d}", r.Url, r.DefaultTTL)
}

type IConfig interface {
	GetAppName() string
	GetAppVersion() string
	GetHttpPort() string
	GetHttpJsonRpcEndpoint() string
	GetHttpMaxRequestBodySize() int64
	GetHttpRequestTimeout() time.Duration
	GetHttpMaxPoolWorkers() int
	GetRedisUrl() string
	GetRedisDefaultTTL() time.Duration
	GetNodes() []NodeConfig
}

type Config struct {
	AppName    string       `yaml:"appName"`
	AppVersion string       `yaml:"appVersion"`
	Http       httpConfig   `yaml:"http"`
	Redis      redisConfig  `yaml:"redis"`
	Nodes      []NodeConfig `yaml:"nodes"`
}

func (c *Config) GetAppName() string {
	return c.AppName
}

func (c *Config) GetAppVersion() string {
	return c.AppVersion
}

func (c *Config) GetHttpPort() string {
	return c.Http.Port
}

func (c *Config) GetHttpJsonRpcEndpoint() string {
	return c.Http.JsonRpcEndpoint
}

func (c *Config) GetHttpMaxRequestBodySize() int64 {
	return c.Http.MaxRequestBodySize
}

func (c *Config) GetHttpRequestTimeout() time.Duration {
	return time.Duration(c.Http.RequestTimeout) * time.Second
}

func (c *Config) GetHttpMaxPoolWorkers() int {
	return c.Http.MaxPoolWorkers
}

func (c *Config) GetRedisUrl() string {
	return c.Redis.Url
}

func (c *Config) GetRedisDefaultTTL() time.Duration {
	return time.Duration(c.Redis.DefaultTTL) * time.Second
}

func (c *Config) GetNodes() []NodeConfig {
	return c.Nodes
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{AppName: %s, AppVersion: %s, Http: %s, Redis: %s, Nodes: %s}",
		c.AppName, c.AppVersion, c.Http.String(), c.Redis.String(), c.Nodes[0].String())
}

var config *Config

func LoadConfig(configPath string, logger logger.ILogger) *Config {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		logger.Error(fmt.Sprintf("Config file does not exist: %s", configPath))
		panic(fmt.Errorf("config file does not exist: %s", configPath))
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read config file: %v", err))
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to unmarshal config YAML: %v", err))
		panic(fmt.Errorf("failed to unmarshal config YAML: %w", err))
	}

	if len(cfg.Nodes) == 0 {
		logger.Panic("Node config is empty")
		panic("node config is empty")
	}

	logger.Log(fmt.Sprintf("Successfully loaded configuration from: %s", configPath))
	return &cfg
}

func SingletonConfig(deps ConfigDeps) *Config {
	if config == nil {
		config = LoadConfig("config.yaml", deps.Logger)
	}
	return config
}
