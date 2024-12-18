package server

import (
	"go.uber.org/dig"
	"json-rpc-node-proxy/internal/handlers"
	"json-rpc-node-proxy/internal/services"
	"json-rpc-node-proxy/pkg/di"
	"json-rpc-node-proxy/pkg/env"
)

func Run(env env.Environment) error {
	c := di.BuildContainer(env)
	c = addAppSpecificDependencies(c)

	return c.Invoke(StartHttpServer)
}

func addAppSpecificDependencies(container *dig.Container) *dig.Container {
	deps := []di.Dependency{
		{
			Constructor: services.NewNodeManagerService,
			Interface:   new(services.INodeManagerService),
			Token:       "NodeManagerService",
		},
		{
			Constructor: services.NewNodeService,
			Interface:   new(services.INodeService),
			Token:       "NodeService",
		},
		{
			Constructor: services.NewJsonRpcService,
			Interface:   new(services.IJsonRpcService),
			Token:       "JsonRpcService",
		},
		{
			Constructor: handlers.NewJsonRpcHandler,
			Interface:   new(handlers.IJsonRpcHandler),
			Token:       "JsonRpcHandler",
		},
	}

	container = di.AppendDependencies(container, deps)

	return container
}
