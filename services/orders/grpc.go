package main

import (
	"log"
	"net"

	handler "github.com/YhomiAce/kitchen/services/orders/handlers/orders"
	"github.com/YhomiAce/kitchen/services/orders/service"
	"google.golang.org/grpc"
)


type gRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}

func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return err
	}
	grpcServer := grpc.NewServer()

	// register grpc services
	orderService := service.NewOrderService()
	handler.NewGrpcOrdersService(grpcServer, orderService)

	log.Println("Starting gRPC server on ", s.addr)

	return grpcServer.Serve(lis)
}