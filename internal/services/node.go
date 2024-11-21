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
	"time"
)

var id = 0

type INodeService interface {
	Request(ctx context.Context) (*models.JsonRpcResponse, error)
}

type NodeService struct {
	logger  logger.ILogger
	url     string
	timeout time.Duration
}

type NodeServiceDependencies struct {
	dig.In

	Logger logger.ILogger     `name:"Logger"`
	Cfg    config.INodeConfig `name:"NodeConfig"`
}

func NewNodeService(deps NodeServiceDependencies) *NodeService {
	return &NodeService{
		logger:  deps.Logger,
		url:     deps.Cfg.GetUrl(),
		timeout: deps.Cfg.GetTimeout(),
	}
}

func (n *NodeService) Request(ctx context.Context) (*models.JsonRpcResponse, error) {
	request, err := utils_ctx.GetJsonRpcRequestFromContext(ctx)

	if err != nil {
		return nil, err
	}

	id += 1

	newRequest := request.CopyWithNewId(fmt.Sprintf("%d", id))

	reqBody, err := json.Marshal(newRequest)

	if err != nil {
		n.logger.Error(fmt.Sprintf("Node.Request() json marshal error %v", err))
		return nil, custom_errors.NodeRequestJsonMarshalError
	}

	respBody, err := n.doRequest(ctx, reqBody)

	if err != nil {
		return nil, err
	}

	jsonRpcResponse := models.NewJsonRpcResponseFromString(string(respBody)).CopyWithNewId(request.Id)

	return jsonRpcResponse, nil
}

func (n *NodeService) doRequest(ctx context.Context, reqBody []byte) ([]byte, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, n.url, bytes.NewReader(reqBody))

	if err != nil {
		n.logger.Error(fmt.Sprintf("Node.doRequest() http.NewRequestWithContext error %v", err))
		return nil, custom_errors.NodeRequestNewRequestError
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: n.timeout}

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
