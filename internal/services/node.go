package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/dig"
	"io/ioutil"
	"json-rpc-node-proxy/internal/models"
	"json-rpc-node-proxy/pkg/config"
	"json-rpc-node-proxy/pkg/custom_errors"
	"json-rpc-node-proxy/pkg/logger"
	utils_ctx "json-rpc-node-proxy/pkg/utils/ctx"
	"net/http"
	"sync"
)

var id = 0

type INodeService interface {
	Request(ctx context.Context) (*models.JsonRpcResponse, error)
}

type NodeService struct {
	manager    INodeManagerService
	logger     logger.ILogger
	reqIdMutex sync.Mutex
}

type NodeServiceDependencies struct {
	dig.In

	NodeManagerService INodeManagerService `name:"NodeManagerService"`
	Logger             logger.ILogger      `name:"Logger"`
	Cfg                config.INodeConfig  `name:"NodeConfig"`
}

func NewNodeService(deps NodeServiceDependencies) *NodeService {
	return &NodeService{
		manager:    deps.NodeManagerService,
		logger:     deps.Logger,
		reqIdMutex: sync.Mutex{},
	}
}

func (n *NodeService) getReqId() int {
	n.reqIdMutex.Lock()
	defer n.reqIdMutex.Unlock()

	id += 1

	return id
}

func (n *NodeService) Request(ctx context.Context) (*models.JsonRpcResponse, error) {
	request, err := utils_ctx.GetJsonRpcRequestFromContext(ctx)

	if err != nil {
		return nil, err
	}

	id := n.getReqId()

	availableNode := n.manager.GetAvailableNodeForRequest(request.Method, id)

	if availableNode == nil {
		n.logger.Error(fmt.Sprintf("Node.Request() no available nodes for method %s", request.Method))
		return nil, custom_errors.AvailableNodeNotFoundError
	}

	newRequest := request.CopyWithNewId(fmt.Sprintf("%d", id))

	reqBody, err := json.Marshal(newRequest)

	if err != nil {
		n.logger.Error(fmt.Sprintf("Node.Request() json marshal error %v", err))
		return nil, custom_errors.NodeRequestJsonMarshalError
	}

	respBody, err := n.doRequest(ctx, availableNode, reqBody)

	if err != nil {
		return nil, err
	}

	jsonRpcResponse := models.NewJsonRpcResponseFromString(string(respBody)).CopyWithNewId(request.Id)

	return jsonRpcResponse, nil
}

func (n *NodeService) doRequest(ctx context.Context, node *models.RateLimitedNode, reqBody []byte) ([]byte, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, node.GetUrl(), bytes.NewReader(reqBody))

	if err != nil {
		n.logger.Error(fmt.Sprintf("Node.doRequest() http.NewRequestWithContext error %v", err))
		return nil, custom_errors.NodeRequestNewRequestError
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: node.GetTimeout()}

	node.WaitForExecute()

	n.logger.Log(fmt.Sprintf("Node.doRequest() sending request to %s", node.GetName()))

	resp, err := client.Do(httpRequest)

	if err != nil {
		n.logger.Error(fmt.Sprintf("Node.doRequest() http client.Do error %v", err))
		return nil, custom_errors.NodeRequestClientDoError
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		n.logger.Error(fmt.Sprintf("Node.doRequest() ioutil.ReadAll error %v", err))
		return nil, custom_errors.NodeRequestReadResponseBodyError
	}

	return respBody, nil
}
