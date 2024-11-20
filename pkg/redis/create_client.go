package redis

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"json-rpc-node-proxy/pkg/config"
)

type CreateRedisClientDependencies struct {
	dig.In
	Cfg config.IRedisConfig `name:"RedisConfig"`
}

func CreateClient(deps CreateRedisClientDependencies) *redis.Client {
	opts, err := redis.ParseURL(deps.Cfg.GetUrl())
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opts)
}
