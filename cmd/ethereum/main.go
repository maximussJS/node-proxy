package main

import (
	"json-rpc-node-proxy/pkg/env"
	"json-rpc-node-proxy/pkg/run_app"
	"log"
)

func main() {
	if err := run_app.RunApp(env.EnvProd); err != nil {
		log.Fatal(err)
	}

	log.Println("Application started")
}
