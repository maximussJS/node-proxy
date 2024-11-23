package config

type INodeConfig interface {
	GetConfigPath() string
}

type NodeConfig struct {
	configPath string
}

func (n *NodeConfig) GetConfigPath() string {
	return n.configPath
}

var nodeConfig *NodeConfig

func init() {
	initNodeConfig()
}

func initNodeConfig() {
	configPath := EnvRequiredString("NODE_CONFIG_PATH")

	nodeConfig = &NodeConfig{
		configPath: configPath,
	}
}

func SingletonNodeConfig() *NodeConfig {
	if nodeConfig == nil {
		initNodeConfig()
	}

	return nodeConfig
}
