package models

import (
	"json-rpc-node-proxy/pkg/algorithms"
	"json-rpc-node-proxy/pkg/config"
	"time"
)

type RateLimitedNode struct {
	node   config.NodeConfig
	Bucket *algorithms.TokenBucket
}

func NewRateLimitedNode(node config.NodeConfig) *RateLimitedNode {
	return &RateLimitedNode{
		node:   node,
		Bucket: algorithms.NewTokenBucket(node.GetTokenBucketCapacity(), node.GetRps()),
	}
}

func (n *RateLimitedNode) GetUrl() string {
	return n.node.GetUrl()
}

func (n *RateLimitedNode) GetName() string {
	return n.node.GetName()
}

func (n *RateLimitedNode) GetTimeout() time.Duration {
	return time.Duration(n.node.GetTimeout()) * time.Second
}

func (n *RateLimitedNode) WaitForExecute() {
	n.Bucket.Wait()
}

func (n *RateLimitedNode) IsWhitelisted(method string) bool {
	return n.node.IsWhitelisted(method)
}

func (n *RateLimitedNode) IsBlacklisted(method string) bool {
	return n.node.IsBlacklisted(method)
}

func (n *RateLimitedNode) IsCached(method string) bool {
	return n.node.IsCached(method)
}
