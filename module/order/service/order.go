package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/trwndh/poc-online-store/module/order/model"
	"github.com/trwndh/poc-online-store/module/order/repository"
	logpkg "github.com/trwndh/poc-online-store/pkg/logger"
)

type service struct {
	ctx       context.Context
	orderRepo repository.OrderRepository
}

type OrderService interface {
	Checkout(ctx context.Context, cart model.Cart) (regPaymentID int64, err error)
	CheckoutNoLocking(ctx context.Context, cart model.Cart) (regPaymentID int64, err error)
}

func NewOrderService(ctx context.Context, orderRepo repository.OrderRepository) OrderService {
	return &service{
		ctx:       ctx,
		orderRepo: orderRepo,
	}
}

var logger *zap.Logger = logpkg.Log
