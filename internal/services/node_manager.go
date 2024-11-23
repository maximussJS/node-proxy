package services

import (
	"go.uber.org/dig"
	"json-rpc-node-proxy/internal/models"
	"json-rpc-node-proxy/pkg/logger"
)

type INodeManagerService interface {
	GetAvailableNodeForRequest(method string, id int) *models.RateLimitedNode
}

type NodeManagerService struct {
	nodes  []*models.RateLimitedNode
	logger logger.ILogger
}

type NodeManagerServiceDependencies struct {
	dig.In

	NodeConfigService INodeConfigService `name:"NodeConfigService"`
	Logger            logger.ILogger     `name:"Logger"`
}

func NewNodeManagerService(deps NodeManagerServiceDependencies) *NodeManagerService {
	nodes := deps.NodeConfigService.GetNodes()

	rateLimitedNodes := make([]*models.RateLimitedNode, len(nodes))

	for i, node := range nodes {
		rateLimitedNodes[i] = models.NewRateLimitedNode(node)
	}

	return &NodeManagerService{
		nodes:  rateLimitedNodes,
		logger: deps.Logger,
	}
}

func (n *NodeManagerService) GetAvailableNodeForRequest(method string, id int) *models.RateLimitedNode {
	availableNodes := n.findAvailableNodeForRequest(method)

	if len(availableNodes) == 0 {
		return nil
	}

	return n.selectNodeRoundRobin(availableNodes, id)
}

func (n *NodeManagerService) findAvailableNodeForRequest(method string) []*models.RateLimitedNode {
	availableNodes := make([]*models.RateLimitedNode, 0)

	for _, node := range n.nodes {
		if node.IsWhitelisted(method) {
			availableNodes = append(availableNodes, node)
			continue
		}

		if node.IsBlacklisted(method) {
			continue
		}

		availableNodes = append(availableNodes, node)
	}

	return availableNodes
}

func (n *NodeManagerService) selectNodeRoundRobin(nodes []*models.RateLimitedNode, reqId int) *models.RateLimitedNode {
	index := reqId % len(nodes)

	return nodes[index]
}

func (n *NodeManagerService) findAvailableNodesForRequest(method string) []*models.RateLimitedNode {
	availableNodes := make([]*models.RateLimitedNode, 0)

	for _, node := range n.nodes {
		if node.IsWhitelisted(method) {
			availableNodes = append(availableNodes, node)
			continue
		}

		if node.IsBlacklisted(method) {
			continue
		}
	}

	return availableNodes
}
