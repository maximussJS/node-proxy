package di

import (
	"go.uber.org/dig"
	"json-rpc-node-proxy/pkg/cache"
	cache_redis "json-rpc-node-proxy/pkg/cache/drivers"
	"json-rpc-node-proxy/pkg/config"
	"json-rpc-node-proxy/pkg/env"
	"json-rpc-node-proxy/pkg/key_generator"
	"json-rpc-node-proxy/pkg/logger"
	redis_cli "json-rpc-node-proxy/pkg/redis"
)

type Dependency struct {
	Constructor interface{}
	Interface   interface{}
	Token       string
}

func BuildContainer(environment env.Environment) *dig.Container {
	dryRun := false
	if environment == env.EnvTest {
		dryRun = true
	}

	deps := getDependencies(environment)

	container := dig.New(dig.DryRun(dryRun))

	for _, dep := range deps {
		mustProvideDependency(container, dep)
	}

	return container
}

func AppendDependencies(container *dig.Container, dependencies []Dependency) *dig.Container {
	for _, dep := range dependencies {
		mustProvideDependency(container, dep)
	}

	return container
}

func getDependencies(env env.Environment) []Dependency {
	return []Dependency{
		{
			Constructor: func() *logger.Logger {
				return logger.NewLogger(env)
			},
			Interface: new(logger.ILogger),
			Token:     "Logger",
		},
		{
			Constructor: config.SingletonConfig,
			Interface:   new(config.IConfig),
			Token:       "Config",
		},
		{
			Constructor: redis_cli.CreateClient,
			Interface:   nil,
			Token:       "RedisClient",
		},
		{
			Constructor: key_generator.NewRedisKeyGenerator,
			Interface:   new(key_generator.IKeyGenerator),
			Token:       "KeyGenerator",
		},
		{
			Constructor: cache_redis.NewRedis,
			Interface:   new(cache.IDriver),
			Token:       "CacheDriver",
		},
		{
			Constructor: cache.NewCache,
			Interface:   new(cache.ICache),
			Token:       "Cache",
		},
	}
}

func mustProvideDependency(container *dig.Container, dependency Dependency) {
	if dependency.Interface == nil {
		err := container.Provide(dependency.Constructor, dig.Name(dependency.Token))

		if err != nil {
			panic(err)
		}

		return
	}

	err := container.Provide(
		dependency.Constructor,
		dig.As(dependency.Interface),
		dig.Name(dependency.Token),
	)

	if err != nil {
		panic(err)
	}
}
