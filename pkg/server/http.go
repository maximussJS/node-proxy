package server

import (
	"github.com/gorilla/mux"
	"go.uber.org/dig"
	"json-rpc-node-proxy/internal/handlers"
	"json-rpc-node-proxy/pkg/config"
	"log"
	"net/http"
)

type HttpServerDependencies struct {
	dig.In
	Cfg            config.IHttpServerConfig `name:"HttpServerConfig"`
	JsonRpcHandler handlers.IJsonRpcHandler `name:"JsonRpcHandler"`
}

func StartHttpServer(deps HttpServerDependencies) error {
	router := mux.NewRouter()

	router.HandleFunc(deps.Cfg.GetJsonRpcEndpoint(), deps.JsonRpcHandler.Handle).Methods("POST")

	port := deps.Cfg.GetPort()

	log.Printf("Starting http server on port %s", port)

	return http.ListenAndServe(port, router)
}
