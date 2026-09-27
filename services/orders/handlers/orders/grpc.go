package handler

import (
	"context"

	"github.com/YhomiAce/kitchen/services/common/genproto/orders"
	"github.com/YhomiAce/kitchen/services/orders/types"
	"google.golang.org/grpc"
)

type OrdersGrpcHandler struct {
	orders.UnimplementedOrderServiceServer

	// service injection
	ordersService types.OrderService
}

func NewGrpcOrdersService(grpc *grpc.Server, ordersService types.OrderService) {
	grpcHandler := &OrdersGrpcHandler{
		ordersService: ordersService,
	}
	
	orders.RegisterOrderServiceServer(grpc, grpcHandler)
}


func(h *OrdersGrpcHandler) CreateOrder(ctx context.Context, request *orders.CreateOrderRequest)(*orders.CreateOrderResponse, error) {
	order := &orders.Order{
		OrderId: 4,
		CustomerId: request.CustomerId,
		ProductId: request.ProductId,
		Quantity: request.Quantity,
	}
	err := h.ordersService.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return &orders.CreateOrderResponse{
		Status: "Success",
	}, nil
}

func(h *OrdersGrpcHandler) GetOrders(ctx context.Context, request *orders.GetOrdersRequest)(*orders.GetOrderResponse, error) {
	ordersData := h.ordersService.GetOrders(ctx, int(request.GetCustomerId()))
	return &orders.GetOrderResponse{
		Orders: ordersData,
	}, nil
}