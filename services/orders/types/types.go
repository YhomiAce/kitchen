package types

import (
	"context"

	"github.com/YhomiAce/kitchen/services/common/genproto/orders"
)


type OrderService interface {
	CreateOrder(context.Context, *orders.Order) error
	GetOrders(context.Context, int) []*orders.Order
}