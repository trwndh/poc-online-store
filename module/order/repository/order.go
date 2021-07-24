package repository

import (
	"context"

	"github.com/trwndh/poc-online-store/module/order/model"
)

type OrderRepository interface {
	Checkout(ctx context.Context, userID int64, cartID int64, cartItems []model.CartItem) (registerPaymentID int64, err error)
	CheckoutNoLocking(ctx context.Context, userID int64, cartID int64, cartItems []model.CartItem) (registerPaymentID int64, err error)
}
