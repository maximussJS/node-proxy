package config

import "fmt"

type NodeConfig struct {
	Name                string   `yaml:"name"`
	Url                 string   `yaml:"url"`
	Timeout             int      `yaml:"timeout"`
	RPS                 int64    `yaml:"rps"`
	TokenBucketCapacity int64    `yaml:"token_bucket_capacity"`
	WhitelistedMethods  []string `yaml:"whitelisted_methods"`
	BlacklistedMethods  []string `yaml:"blacklisted_methods"`
	CachedMethods       []string `yaml:"cached_methods"`
}

func (n *NodeConfig) GetName() string {
	return n.Name
}

func (n *NodeConfig) GetUrl() string {
	return n.Url
}

func (n *NodeConfig) GetTimeout() int {
	return n.Timeout
}

func (n *NodeConfig) GetRps() int64 {
	return n.RPS
}

func (n *NodeConfig) String() string {
	return fmt.Sprintf("NodeConfig{name=%s, url=%s, timeout=%d, rps=%d, tokenBucketCapacity=%d, whitelistedMethods=%v, blacklistedMethods=%v, cachedMethods=%v}",
		n.Name, n.Url, n.Timeout, n.RPS, n.TokenBucketCapacity, n.WhitelistedMethods, n.BlacklistedMethods, n.CachedMethods)
}

func (n *NodeConfig) GetTokenBucketCapacity() int64 {
	return n.TokenBucketCapacity
}

func (n *NodeConfig) IsWhitelisted(method string) bool {
	for _, m := range n.WhitelistedMethods {
		if m == method {
			return true
		}
	}

	return false
}

func (n *NodeConfig) IsBlacklisted(method string) bool {
	for _, m := range n.BlacklistedMethods {
		if m == method {
			return true
		}
	}

	return false
}

func (n *NodeConfig) IsCached(method string) bool {
	for _, m := range n.CachedMethods {
		if m == method {
			return true
		}
	}

	return false
}
