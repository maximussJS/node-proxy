package node

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"json-rpc-node-proxy/pkg/node-proxy/common/custom_errors"
	"json-rpc-node-proxy/pkg/node-proxy/config"
	"json-rpc-node-proxy/pkg/node-proxy/json-rpc-types"
	"log"
	"net/http"
	"time"
)

var id = 0

type Node struct {
	url     string
	timeout time.Duration
}

func NewNode() *Node {
	return &Node{
		url:     config.NodeConfig.Url,
		timeout: time.Duration(config.NodeConfig.Timeout) * time.Second,
	}
}

func (n *Node) Request(ctx context.Context, r *json_rpc_types.JsonRpcRequest) (*json_rpc_types.JsonRpcResponse, error) {
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

	jsonRpcResponse := json_rpc_types.NewJsonRpcResponseFromString(string(respBody)).CopyWithNewId(r.Id)

	return jsonRpcResponse, nil
}
