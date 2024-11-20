package main

import (
	"json-rpc-node-proxy/cmd/runner"
	"json-rpc-node-proxy/pkg/env"
	"log"
)

func main() {
	if err := runner.RunApp(env.EnvProd); err != nil {
		log.Fatal(err)
	}

	log.Println("Application started")
}
