package redis

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"json-rpc-node-proxy/pkg/config"
)

type CreateRedisClientDependencies struct {
	dig.In
	Cfg config.IConfig `name:"Config"`
}

func CreateClient(deps CreateRedisClientDependencies) *redis.Client {
	opts, err := redis.ParseURL(deps.Cfg.GetRedisUrl())
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opts)
}
