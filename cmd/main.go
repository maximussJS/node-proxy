package main

import (
	"json-rpc-node-proxy/pkg/env"
	"json-rpc-node-proxy/pkg/server"
	"log"
)

func main() {
	if err := server.Run(env.EnvProd); err != nil {
		log.Fatal(err)
	}
}
