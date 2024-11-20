package config

import "time"

type INodeConfig interface {
	GetUrl() string
	GetTimeout() time.Duration
}

type NodeConfig struct {
	Url     string
	Timeout int
}

func (nc *NodeConfig) GetUrl() string {
	return nc.Url
}

func (nc *NodeConfig) GetTimeout() time.Duration {
	return time.Duration(nc.Timeout) * time.Second
}

var nodeConfig *NodeConfig

func init() {
	initNodeConfig()
}

func initNodeConfig() {
	url := EnvRequiredString("NODE_URL")

	timeout := EnvOptionalInt("NODE_TIMEOUT", 30)

	nodeConfig = &NodeConfig{
		Url:     url,
		Timeout: timeout,
	}
}

func SingletonNodeConfig() *NodeConfig {
	if nodeConfig == nil {
		initNodeConfig()
	}

	return nodeConfig
}
