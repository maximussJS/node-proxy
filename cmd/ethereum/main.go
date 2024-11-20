package main

import (
	"go.uber.org/dig"
	"json-rpc-node-proxy/internal/handlers"
	"json-rpc-node-proxy/internal/services"
	"json-rpc-node-proxy/pkg/di"
	"json-rpc-node-proxy/pkg/server"
	"log"
)

func main() {
	if err := runApp(di.EnvProd); err != nil {
		log.Fatal(err)
	}
}

func runApp(env di.Environment) error {
	c := di.BuildContainer(env)
	c = addAppSpecificDependencies(c)

	return c.Invoke(server.StartHttpServer)
}

func addAppSpecificDependencies(container *dig.Container) *dig.Container {
	deps := []di.Dependency{
		{
			Constructor: services.NewNodeService,
			Interface:   new(services.INodeService),
			Token:       "NodeService",
		},
		{
			Constructor: services.NewProxyService,
			Interface:   new(services.IProxyService),
			Token:       "ProxyService",
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
