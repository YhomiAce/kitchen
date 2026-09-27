package handler

import (
	"net/http"

	"github.com/YhomiAce/kitchen/services/common/genproto/orders"
	"github.com/YhomiAce/kitchen/services/common/util"
	"github.com/YhomiAce/kitchen/services/orders/types"
)

type OrderHttpHandler struct {
	orderService types.OrderService
}

func NewOrderHttpHandler(orderService types.OrderService) *OrderHttpHandler {
	return &OrderHttpHandler{orderService: orderService}
}

func(h *OrderHttpHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /orders", h.CreateOrder)
}

func(h *OrderHttpHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req orders.CreateOrderRequest
	err := util.ParseJSON(r, &req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
	}

	order := &orders.Order{
		OrderId: 5,
		CustomerId: req.CustomerId,
		ProductId: req.ProductId,
		Quantity: req.Quantity,
	}

	err = h.orderService.CreateOrder(r.Context(), order)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
	res := &orders.CreateOrderResponse{Status: "success"}
	util.WriteJSON(w, http.StatusCreated, res)
}