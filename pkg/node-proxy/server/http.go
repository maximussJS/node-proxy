package server

import (
	"github.com/gorilla/mux"
	"json-rpc-node-proxy/pkg/node-proxy/config"
	rpc "json-rpc-node-proxy/pkg/node-proxy/json-rpc"
	"log"
	"net/http"
	"time"
)

type HttpServer struct {
	router *mux.Router
}

func NewHttpServer() *HttpServer {
	server := &HttpServer{
		router: mux.NewRouter(),
	}

	server.initRoutes()

	return server
}

func (s *HttpServer) initRoutes() {
	jsonRpcHandler := rpc.NewJsonRpcHandler(
		config.HttpServerConfig.MaxRequestBodySize,
		time.Duration(config.HttpServerConfig.Timeout)*time.Second,
	)

	s.router.HandleFunc(config.HttpServerConfig.JsonRpcEndpoint, jsonRpcHandler.Handle).Methods("POST")
}

func (s *HttpServer) ListenAndServe() {
	port := config.HttpServerConfig.Port

	log.Printf("Starting http server on port %s", port)

	panic(http.ListenAndServe(port, s.router))
}
