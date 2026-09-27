package main

import (
	"log"
	"net/http"

	handler "github.com/YhomiAce/kitchen/services/orders/handlers/orders"
	"github.com/YhomiAce/kitchen/services/orders/service"
)

type httpServer struct {
	addr string
}

func NewHttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()

	orderService := service.NewOrderService()
	orderHandler := handler.NewOrderHttpHandler(orderService)
	orderHandler.RegisterRoutes(router)

	log.Println("Starting server on: ", s.addr)
	return http.ListenAndServe(s.addr, router)
}