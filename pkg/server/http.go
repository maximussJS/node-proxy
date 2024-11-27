package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/dig"
	"json-rpc-node-proxy/internal/handlers"
	"json-rpc-node-proxy/pkg/config"
	"json-rpc-node-proxy/pkg/logger"
	"net/http"
)

type HttpServerDependencies struct {
	dig.In

	Logger         logger.ILogger           `name:"Logger"`
	Cfg            config.IConfig           `name:"Config"`
	JsonRpcHandler handlers.IJsonRpcHandler `name:"JsonRpcHandler"`
}

func StartHttpServer(deps HttpServerDependencies) error {
	router := mux.NewRouter()

	router.HandleFunc(deps.Cfg.GetHttpJsonRpcEndpoint(), deps.JsonRpcHandler.Handle).Methods("POST")

	port := deps.Cfg.GetHttpPort()

	deps.Logger.Log(fmt.Sprintf("Starting http server on port %s", port))

	return http.ListenAndServe(port, router)
}
