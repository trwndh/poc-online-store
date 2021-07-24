package handler

import (
	"net/http"

	order "github.com/trwndh/poc-online-store/module/order/service"
)

type handler struct {
	orderService order.OrderService
}

type Handler interface {
	CheckoutHandler(w http.ResponseWriter, r *http.Request)
}

func NewHandler(orderService order.OrderService) Handler {
	return &handler{orderService: orderService}
}
