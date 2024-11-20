package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/dig"
	"io/ioutil"
	"json-rpc-node-proxy/internal/common/custom_errors"
	"json-rpc-node-proxy/internal/models"
	"json-rpc-node-proxy/pkg/config"
	"log"
	"net/http"
	"time"
)

var id = 0

type INodeService interface {
	Request(ctx context.Context, r *models.JsonRpcRequest) (*models.JsonRpcResponse, error)
}

type NodeService struct {
	url     string
	timeout time.Duration
}

type NodeServiceDependencies struct {
	dig.In
	Cfg config.INodeConfig `name:"NodeConfig"`
}

func NewNodeService(deps NodeServiceDependencies) *NodeService {
	return &NodeService{
		url:     deps.Cfg.GetUrl(),
		timeout: deps.Cfg.GetTimeout(),
	}
}

func (n *NodeService) Request(ctx context.Context, r *models.JsonRpcRequest) (*models.JsonRpcResponse, error) {
	id += 1

	newRequest := r.CopyWithNewId(fmt.Sprintf("%d", id))

	reqBody, err := json.Marshal(newRequest)
	if err != nil {
		log.Printf("Node.Request() json marshal error %v", err)
		return nil, custom_errors.NodeRequestJsonMarshalError
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, n.url, bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("Node.Request() http.NewRequestWithContext error %v", err)
		return nil, custom_errors.NodeRequestNewRequestError
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: n.timeout}
	resp, err := client.Do(httpRequest)
	if err != nil {
		log.Printf("Node.Request() http client.Do error %v", err)
		return nil, custom_errors.NodeRequestClientDoError
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Node.Request() ioutil.ReadAll error %v", err)
		return nil, custom_errors.NodeRequestReadResponseBodyError
	}

	jsonRpcResponse := models.NewJsonRpcResponseFromString(string(respBody)).CopyWithNewId(r.Id)

	return jsonRpcResponse, nil
}
