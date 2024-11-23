package models

type Node struct {
	Name                string   `yaml:"name"`
	Url                 string   `yaml:"url"`
	Timeout             int64    `yaml:"timeout"`
	Rps                 int64    `yaml:"rps"`
	TokenBucketCapacity int64    `yaml:"token_bucket_capacity"`
	WhitelistedMethods  []string `yaml:"whitelisted_methods"`
	BlacklistedMethods  []string `yaml:"blacklisted_methods"`
	CachedMethods       []string `yaml:"cached_methods"`
}

func (n *Node) IsWhitelisted(method string) bool {
	for _, m := range n.WhitelistedMethods {
		if m == method {
			return true
		}
	}

	return false
}

func (n *Node) IsBlacklisted(method string) bool {
	for _, m := range n.BlacklistedMethods {
		if m == method {
			return true
		}
	}

	return false
}

func (n *Node) IsCached(method string) bool {
	for _, m := range n.CachedMethods {
		if m == method {
			return true
		}
	}

	return false
}
