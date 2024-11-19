package main

import (
	"json-rpc-node-proxy/pkg/node-proxy/server"
)

func main() {

	httpServer := server.NewHttpServer()

	httpServer.ListenAndServe()
}
