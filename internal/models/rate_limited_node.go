package models

import (
	"json-rpc-node-proxy/pkg/algorithms"
	"time"
)

type RateLimitedNode struct {
	node   Node
	Bucket *algorithms.TokenBucket
}

func NewRateLimitedNode(node Node) *RateLimitedNode {
	return &RateLimitedNode{
		node:   node,
		Bucket: algorithms.NewTokenBucket(node.TokenBucketCapacity, node.Rps),
	}
}

func (n *RateLimitedNode) GetUrl() string {
	return n.node.Url
}

func (n *RateLimitedNode) GetName() string {
	return n.node.Name
}

func (n *RateLimitedNode) GetTimeout() time.Duration {
	return time.Duration(n.node.Timeout) * time.Second
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