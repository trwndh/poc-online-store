package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/trwndh/poc-online-store/module/product/model"
)

func (s *service) GetProductByID(ctx context.Context, productID int64) (product model.Product, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.product.GetProductByID")
	defer func() {
		if err != nil {
			span.SetTag("error", true).LogFields(log.String("error", err.Error()))
		}
		span.Finish()
	}()

	product, err = s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		logger.Error("error when get products: " + err.Error())
		return product, err
	}

	return product, nil
}
