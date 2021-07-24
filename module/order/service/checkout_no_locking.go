package service

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/trwndh/poc-online-store/module/order/model"
	"go.uber.org/zap"
)

func (s *service) CheckoutNoLocking(ctx context.Context, cart model.Cart) (regPaymentID int64, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.order.checkout")
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
		}
		span.Finish()
	}()

	if cart.IsEmptyItems() || cart.IsUserIDNotValid() {
		msg := "empty item on the cart"
		if !cart.IsUserIDNotValid() {
			msg = "invalid user id"
		}
		logger.Error("Invalid input given: " + msg)
		return 0, err
	}

	// unmarshal items to cart item struct
	regPaymentID, err = s.orderRepo.CheckoutNoLocking(ctx, cart.UserID, cart.ID, cart.Items)
	if err != nil {
		logger.Error(fmt.Sprintf("error when checkout: %s ", err.Error()), zap.Any("userID", cart.UserID), zap.Any("cartID", cart.ID))
		return 0, err
	}

	return regPaymentID, nil
}
