package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/YhomiAce/kitchen/services/common/genproto/orders"
	"github.com/YhomiAce/kitchen/services/common/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type httpServer struct {
	addr string
}

func NewGRPCClient(addr string)*grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Could not connect to GRPC server", err)
		return nil
	}
	log.Println("Connected to GRPC Server on ", addr)
	return conn
}

func NewHttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()

	conn := NewGRPCClient(":9000")
	defer conn.Close()

	client := orders.NewOrderServiceClient(conn)

	router.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		var req orders.CreateOrderRequest
		err := util.ParseJSON(r, &req)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err)
			return 
		}

		res, err := client.CreateOrder(r.Context(), &req)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err)
			return ;
		}
		util.WriteJSON(w, http.StatusCreated, res)
	})

	router.HandleFunc("GET /orders/{customerId}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("customerId")
		customerId, _ := strconv.Atoi(id)
		req := &orders.GetOrdersRequest{
			CustomerId: int32(customerId),
		}
		res, err := client.GetOrders(r.Context(), req)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err)
			return 
		}
		util.WriteJSON(w, http.StatusOK, res)
	})
	
	log.Println("Starting server on: ", s.addr)
	return http.ListenAndServe(s.addr, router)
}