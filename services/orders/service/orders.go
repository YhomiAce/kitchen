package service

import (
	"context"

	"github.com/YhomiAce/kitchen/services/common/genproto/orders"
)

var ordersDb = make([] *orders.Order, 0)


type OrderService struct {
	
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func(s *OrderService)CreateOrder(ctx context.Context, order *orders.Order) error {
	ordersDb = append(ordersDb, order)
	return nil
}

func(s *OrderService) GetOrders(ctx context.Context, customerId int) []*orders.Order {
	return ordersDb
}