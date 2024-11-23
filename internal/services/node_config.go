package services

import (
	"errors"
	"fmt"
	"go.uber.org/dig"
	"gopkg.in/yaml.v3"
	"json-rpc-node-proxy/internal/models"
	"json-rpc-node-proxy/pkg/config"
	"json-rpc-node-proxy/pkg/logger"
	"os"
)

type INodeConfigService interface {
	GetNodes() []models.Node
}

type NodeConfigServiceDependencies struct {
	dig.In

	Logger logger.ILogger     `name:"Logger"`
	Cfg    config.INodeConfig `name:"NodeConfig"`
}

type NodeConfigService struct {
	nodes []models.Node
}

func NewNodeConfigService(deps NodeConfigServiceDependencies) *NodeConfigService {
	configPath := deps.Cfg.GetConfigPath()

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		deps.Logger.Error(fmt.Sprintf("Config file does not exist: %s", configPath))
		panic(fmt.Errorf("config file does not exist: %s", configPath))
	}

	nodes, err := parseNodeConfig(configPath)
	if err != nil {
		deps.Logger.Error(fmt.Sprintf("Failed to parse node config: %v", err))
		panic(fmt.Errorf("failed to parse node config: %w", err))
	}

	deps.Logger.Log(fmt.Sprintf("Successfully parsed node config: %v", configPath))

	if len(nodes) == 0 {
		deps.Logger.Panic("Node config is empty")
		panic("node config is empty")
	}

	return &NodeConfigService{
		nodes: nodes,
	}
}

func (n *NodeConfigService) GetNodes() []models.Node {
	return n.nodes
}

func parseNodeConfig(configPath string) ([]models.Node, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var nodes []models.Node
	err = yaml.Unmarshal(data, &nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return nodes, nil
}
